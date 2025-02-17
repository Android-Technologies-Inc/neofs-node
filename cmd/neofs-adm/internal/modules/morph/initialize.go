package morph

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/nspcc-dev/neo-go/pkg/core/native/nativenames"
	"github.com/nspcc-dev/neo-go/pkg/core/state"
	"github.com/nspcc-dev/neo-go/pkg/core/transaction"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/rpc/client"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract/trigger"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"github.com/nspcc-dev/neo-go/pkg/wallet"
	"github.com/nspcc-dev/neofs-node/cmd/neofs-adm/internal/modules/config"
	"github.com/nspcc-dev/neofs-node/pkg/innerring"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type cache struct {
	nnsCs    *state.Contract
	groupKey *keys.PublicKey
}

type initializeContext struct {
	clientContext
	cache
	// CommitteeAcc is used for retrieving committee address and verification script.
	CommitteeAcc *wallet.Account
	// ConsensusAcc is used for retrieving committee address and verification script.
	ConsensusAcc *wallet.Account
	Wallets      []*wallet.Wallet
	// ContractWallet is a wallet for providing contract group signature.
	ContractWallet *wallet.Wallet
	// Accounts contains simple signature accounts in the same order as in Wallets.
	Accounts     []*wallet.Account
	Contracts    map[string]*contractState
	Command      *cobra.Command
	ContractPath string
	Natives      map[string]util.Uint160
}

func initializeSideChainCmd(cmd *cobra.Command, args []string) error {
	initCtx, err := newInitializeContext(cmd, viper.GetViper())
	if err != nil {
		return fmt.Errorf("initialization error: %w", err)
	}

	// 1. Transfer funds to committee accounts.
	cmd.Println("Stage 1: transfer GAS to alphabet nodes.")
	if err := initCtx.transferFunds(); err != nil {
		return err
	}

	cmd.Println("Stage 2: set notary and alphabet nodes in designate contract.")
	if err := initCtx.setNotaryAndAlphabetNodes(); err != nil {
		return err
	}

	// 3. Deploy NNS contract.
	cmd.Println("Stage 3: deploy NNS contract.")
	if err := initCtx.deployNNS(deployMethodName); err != nil {
		return err
	}

	// 4. Deploy NeoFS contracts.
	cmd.Println("Stage 4: deploy NeoFS contracts.")
	if err := initCtx.deployContracts(); err != nil {
		return err
	}

	cmd.Println("Stage 4.1: Transfer GAS to proxy contract.")
	if err := initCtx.transferGASToProxy(); err != nil {
		return err
	}

	cmd.Println("Stage 5: register candidates.")
	if err := initCtx.registerCandidates(); err != nil {
		return err
	}

	cmd.Println("Stage 6: transfer NEO to alphabet contracts.")
	if err := initCtx.transferNEOToAlphabetContracts(); err != nil {
		return err
	}

	cmd.Println("Stage 7: set addresses in NNS.")
	if err := initCtx.setNNS(); err != nil {
		return err
	}

	return nil
}

func newInitializeContext(cmd *cobra.Command, v *viper.Viper) (*initializeContext, error) {
	walletDir := config.ResolveHomePath(viper.GetString(alphabetWalletsFlag))
	wallets, err := openAlphabetWallets(walletDir)
	if err != nil {
		return nil, err
	}

	w, err := openContractWallet(cmd, walletDir)
	if err != nil {
		return nil, err
	}

	c, err := getN3Client(v)
	if err != nil {
		return nil, fmt.Errorf("can't create N3 client: %w", err)
	}

	committeeAcc, err := getWalletAccount(wallets[0], committeeAccountName)
	if err != nil {
		return nil, fmt.Errorf("can't find committee account: %w", err)
	}

	consensusAcc, err := getWalletAccount(wallets[0], consensusAccountName)
	if err != nil {
		return nil, fmt.Errorf("can't find consensus account: %w", err)
	}

	var ctrPath string
	if cmd.Name() == "init" {
		if viper.GetInt64(epochDurationInitFlag) <= 0 {
			return nil, fmt.Errorf("epoch duration must be positive")
		}

		if viper.GetInt64(maxObjectSizeInitFlag) <= 0 {
			return nil, fmt.Errorf("max object size must be positive")
		}
	}

	needContracts := cmd.Name() == "update-contracts" || cmd.Name() == "init"
	if needContracts {
		ctrPath, err = cmd.Flags().GetString(contractsInitFlag)
		if err != nil {
			return nil, fmt.Errorf("invalid contracts path: %w", err)
		}
	}

	nativeHashes, err := getNativeHashes(c)
	if err != nil {
		return nil, err
	}

	accounts := make([]*wallet.Account, len(wallets))
	for i, w := range wallets {
		acc, err := getWalletAccount(w, singleAccountName)
		if err != nil {
			return nil, fmt.Errorf("wallet %s is invalid (no single account): %w", w.Path(), err)
		}
		accounts[i] = acc
	}

	initCtx := &initializeContext{
		clientContext:  *defaultClientContext(c),
		ConsensusAcc:   consensusAcc,
		CommitteeAcc:   committeeAcc,
		ContractWallet: w,
		Wallets:        wallets,
		Accounts:       accounts,
		Command:        cmd,
		Contracts:      make(map[string]*contractState),
		ContractPath:   ctrPath,
		Natives:        nativeHashes,
	}

	if needContracts {
		err := initCtx.readContracts(fullContractList)
		if err != nil {
			return nil, err
		}
	}

	return initCtx, nil
}

func (c *initializeContext) nativeHash(name string) util.Uint160 {
	return c.Natives[name]
}

func openAlphabetWallets(walletDir string) ([]*wallet.Wallet, error) {
	walletFiles, err := ioutil.ReadDir(walletDir)
	if err != nil {
		return nil, fmt.Errorf("can't read alphabet wallets dir: %w", err)
	}

	var size int
loop:
	for i := range walletFiles {
		for j := 0; j < len(walletFiles); j++ {
			letter := innerring.GlagoliticLetter(j).String()
			if strings.HasPrefix(walletFiles[i].Name(), letter) {
				size++
				continue loop
			}
		}
	}
	if size == 0 {
		return nil, errors.New("alphabet wallets dir is empty (run `generate-alphabet` command first)")
	}

	wallets := make([]*wallet.Wallet, size)
	for i := 0; i < size; i++ {
		p := filepath.Join(walletDir, innerring.GlagoliticLetter(i).String()+".json")
		w, err := wallet.NewWalletFromFile(p)
		if err != nil {
			return nil, fmt.Errorf("can't open wallet: %w", err)
		}

		password, err := config.AlphabetPassword(viper.GetViper(), i)
		if err != nil {
			return nil, fmt.Errorf("can't fetch password: %w", err)
		}

		for i := range w.Accounts {
			if err := w.Accounts[i].Decrypt(password, keys.NEP2ScryptParams()); err != nil {
				return nil, fmt.Errorf("can't unlock wallet: %w", err)
			}
		}

		wallets[i] = w
	}

	return wallets, nil
}

func (c *initializeContext) awaitTx() error {
	return c.clientContext.awaitTx(c.Command)
}

func (c *initializeContext) nnsContractState() (*state.Contract, error) {
	if c.nnsCs != nil {
		return c.nnsCs, nil
	}

	cs, err := c.Client.GetContractStateByID(1)
	if err != nil {
		return nil, err
	}

	c.nnsCs = cs
	return cs, nil
}

func (c *initializeContext) getSigner() transaction.Signer {
	if c.groupKey != nil {
		return transaction.Signer{
			Scopes:        transaction.CustomGroups,
			AllowedGroups: keys.PublicKeys{c.groupKey},
		}
	}

	signer := transaction.Signer{
		Account: c.CommitteeAcc.Contract.ScriptHash(),
		Scopes:  transaction.Global, // Scope is important, as we have nested call to container contract.
	}

	nnsCs, err := c.nnsContractState()
	if err != nil {
		return signer
	}

	groupKey, err := nnsResolveKey(c.Client, nnsCs.Hash, groupKeyDomain)
	if err == nil {
		c.groupKey = groupKey

		signer.Scopes = transaction.CustomGroups
		signer.AllowedGroups = keys.PublicKeys{groupKey}
	}
	return signer
}

func (c *clientContext) awaitTx(cmd *cobra.Command) error {
	if len(c.Hashes) == 0 {
		return nil
	}

	cmd.Println("Waiting for transactions to persist...")

	tick := time.NewTicker(c.PollInterval)
	defer tick.Stop()

	timer := time.NewTimer(c.WaitDuration)
	defer timer.Stop()

	at := trigger.Application

loop:
	for i := range c.Hashes {
		_, err := c.Client.GetApplicationLog(c.Hashes[i], &at)
		if err == nil {
			continue loop
		}
		for {
			select {
			case <-tick.C:
				_, err := c.Client.GetApplicationLog(c.Hashes[i], &at)
				if err == nil {
					continue loop
				}
			case <-timer.C:
				return errors.New("timeout while waiting for transaction to persist")
			}
		}
	}

	return nil
}

func (c *initializeContext) sendCommitteeTx(script []byte, sysFee int64) error {
	tx, err := c.Client.CreateTxFromScript(script, c.CommitteeAcc, sysFee, 0, []client.SignerAccount{{
		Signer:  c.getSigner(),
		Account: c.CommitteeAcc,
	}})
	if err != nil {
		return fmt.Errorf("can't create tx: %w", err)
	}

	return c.multiSignAndSend(tx, committeeAccountName)
}

func getWalletAccount(w *wallet.Wallet, typ string) (*wallet.Account, error) {
	for i := range w.Accounts {
		if w.Accounts[i].Label == typ {
			return w.Accounts[i], nil
		}
	}
	return nil, fmt.Errorf("account for '%s' not found", typ)
}

func getNativeHashes(c *client.Client) (map[string]util.Uint160, error) {
	ns, err := c.GetNativeContracts()
	if err != nil {
		return nil, fmt.Errorf("can't get native contract hashes: %w", err)
	}

	notaryEnabled := false
	nativeHashes := make(map[string]util.Uint160, len(ns))
	for i := range ns {
		if ns[i].Manifest.Name == nativenames.Notary {
			notaryEnabled = len(ns[i].UpdateHistory) > 0
		}
		nativeHashes[ns[i].Manifest.Name] = ns[i].Hash
	}
	if !notaryEnabled {
		return nil, errors.New("notary contract must be enabled")
	}
	return nativeHashes, nil
}
