package solanaclient

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

const (
	PROGRAM_ID           = "6xYHZd8RXahNYfUxVRGSLk22n97gMynGGEpThNCQKwaB"
	CASH_ACCOUNT_SEED    = "cash-account"
	PENDING_REQUEST_SEED = "pending-request"
)

type CashDappClient struct {
	client    *rpc.Client
	programID solana.PublicKey
}

func NewCashDappClient(endpoint string) (*CashDappClient, error) {
	programID, err := solana.PublicKeyFromBase58(PROGRAM_ID)
	if err != nil {
		return nil, fmt.Errorf("invalid program ID: %v", err)
	}

	return &CashDappClient{
		client:    rpc.New(endpoint),
		programID: programID,
	}, nil
}

// InitializeAccount creates a new cash account for the given wallet
func (c *CashDappClient) InitializeAccount(ctx context.Context, wallet *solana.Wallet) (string, error) {
	// Derive PDA for cash account
	cashAccountPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(CASH_ACCOUNT_SEED),
			wallet.PublicKey().Bytes(),
		},
		c.programID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to derive cash account PDA: %v", err)
	}

	// Get recent blockhash
	recent, err := c.client.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	// Create instruction data
	data := []byte{74, 115, 99, 93, 197, 69, 103, 7} // initialize_account discriminator

	// Create instruction
	ix := solana.NewInstruction(
		c.programID,
		solana.AccountMetaSlice{
			{PublicKey: cashAccountPDA, IsSigner: false, IsWritable: true},
			{PublicKey: wallet.PublicKey(), IsSigner: true, IsWritable: true},
			{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
		},
		data,
	)

	// Create and sign transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(wallet.PublicKey()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %v", err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if wallet.PublicKey().Equals(key) {
			return &wallet.PrivateKey
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send transaction
	sig, err := c.client.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return sig.String(), nil
}

// DepositFunds deposits SOL into the cash account
func (c *CashDappClient) DepositFunds(ctx context.Context, wallet *solana.Wallet, amount uint64) (string, error) {
	// Derive PDA for cash account
	cashAccountPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(CASH_ACCOUNT_SEED),
			wallet.PublicKey().Bytes(),
		},
		c.programID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to derive cash account PDA: %v", err)
	}

	// Get recent blockhash
	recent, err := c.client.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	// Create instruction data
	data := append(
		[]byte{202, 39, 52, 211, 53, 20, 250, 88}, // deposit_funds discriminator
		binary.BigEndian.AppendUint64(nil, amount)...,
	)

	// Create instruction
	ix := solana.NewInstruction(
		c.programID,
		solana.AccountMetaSlice{
			{PublicKey: cashAccountPDA, IsSigner: false, IsWritable: true},
			{PublicKey: wallet.PublicKey(), IsSigner: true, IsWritable: true},
			{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
		},
		data,
	)

	// Create and sign transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(wallet.PublicKey()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %v", err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if wallet.PublicKey().Equals(key) {
			return &wallet.PrivateKey
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send transaction
	sig, err := c.client.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return sig.String(), nil
}

// TransferFunds transfers SOL from one cash account to another
func (c *CashDappClient) TransferFunds(ctx context.Context, wallet *solana.Wallet, recipient solana.PublicKey, amount uint64) (string, error) {
	// Derive PDA for sender's cash account
	fromCashAccountPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(CASH_ACCOUNT_SEED),
			wallet.PublicKey().Bytes(),
		},
		c.programID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to derive sender cash account PDA: %v", err)
	}

	// Derive PDA for recipient's cash account
	toCashAccountPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(CASH_ACCOUNT_SEED),
			recipient.Bytes(),
		},
		c.programID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to derive recipient cash account PDA: %v", err)
	}

	// Get recent blockhash
	recent, err := c.client.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	// Create instruction data
	data := append(
		[]byte{238, 99, 187, 231, 93, 119, 76, 238}, // transfer_funds discriminator
		append(
			recipient.Bytes(),
			binary.BigEndian.AppendUint64(nil, amount)...,
		)...,
	)

	// Create instruction
	ix := solana.NewInstruction(
		c.programID,
		solana.AccountMetaSlice{
			{PublicKey: fromCashAccountPDA, IsSigner: false, IsWritable: true},
			{PublicKey: toCashAccountPDA, IsSigner: false, IsWritable: true},
			{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
			{PublicKey: wallet.PublicKey(), IsSigner: true, IsWritable: false},
		},
		data,
	)

	// Create and sign transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(wallet.PublicKey()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %v", err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if wallet.PublicKey().Equals(key) {
			return &wallet.PrivateKey
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send transaction
	sig, err := c.client.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return sig.String(), nil
}

// WithdrawFunds withdraws SOL from the cash account
func (c *CashDappClient) WithdrawFunds(ctx context.Context, wallet *solana.Wallet, amount uint64) (string, error) {
	// Derive PDA for cash account
	cashAccountPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(CASH_ACCOUNT_SEED),
			wallet.PublicKey().Bytes(),
		},
		c.programID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to derive cash account PDA: %v", err)
	}

	// Get recent blockhash
	recent, err := c.client.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	// Create instruction data
	data := append(
		[]byte{241, 36, 29, 111, 208, 31, 104, 217}, // withdraw_funds discriminator
		binary.BigEndian.AppendUint64(nil, amount)...,
	)

	// Create instruction
	ix := solana.NewInstruction(
		c.programID,
		solana.AccountMetaSlice{
			{PublicKey: cashAccountPDA, IsSigner: false, IsWritable: true},
			{PublicKey: wallet.PublicKey(), IsSigner: true, IsWritable: true},
			{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
		},
		data,
	)

	// Create and sign transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(wallet.PublicKey()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %v", err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if wallet.PublicKey().Equals(key) {
			return &wallet.PrivateKey
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send transaction
	sig, err := c.client.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return sig.String(), nil
}

// AddFriend adds a friend to the cash account
func (c *CashDappClient) AddFriend(ctx context.Context, wallet *solana.Wallet, friendPubKey solana.PublicKey) (string, error) {
	// Derive PDA for cash account
	cashAccountPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(CASH_ACCOUNT_SEED),
			wallet.PublicKey().Bytes(),
		},
		c.programID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to derive cash account PDA: %v", err)
	}

	// Get recent blockhash
	recent, err := c.client.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	// Create instruction data
	data := append(
		[]byte{6, 45, 26, 157, 246, 216, 236, 32}, // add_friend discriminator
		friendPubKey.Bytes()...,
	)

	// Create instruction
	ix := solana.NewInstruction(
		c.programID,
		solana.AccountMetaSlice{
			{PublicKey: cashAccountPDA, IsSigner: false, IsWritable: true},
			{PublicKey: wallet.PublicKey(), IsSigner: true, IsWritable: true},
			{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
		},
		data,
	)

	// Create and sign transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(wallet.PublicKey()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %v", err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if wallet.PublicKey().Equals(key) {
			return &wallet.PrivateKey
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send transaction
	sig, err := c.client.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return sig.String(), nil
}

// NewPendingRequest creates a new pending request for funds
func (c *CashDappClient) NewPendingRequest(ctx context.Context, wallet *solana.Wallet, sender solana.PublicKey, amount uint64) (string, error) {
	// Derive PDA for pending request
	pendingRequestPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(PENDING_REQUEST_SEED),
			wallet.PublicKey().Bytes(),
		},
		c.programID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to derive pending request PDA: %v", err)
	}

	// Derive PDA for cash account
	cashAccountPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(CASH_ACCOUNT_SEED),
			wallet.PublicKey().Bytes(),
		},
		c.programID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to derive cash account PDA: %v", err)
	}

	// Get recent blockhash
	recent, err := c.client.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	// Create instruction data
	data := append(
		[]byte{107, 236, 75, 202, 60, 9, 222, 51}, // new_pending_request discriminator
		append(
			sender.Bytes(),
			binary.BigEndian.AppendUint64(nil, amount)...,
		)...,
	)

	// Create instruction
	ix := solana.NewInstruction(
		c.programID,
		solana.AccountMetaSlice{
			{PublicKey: pendingRequestPDA, IsSigner: false, IsWritable: true},
			{PublicKey: cashAccountPDA, IsSigner: false, IsWritable: true},
			{PublicKey: wallet.PublicKey(), IsSigner: true, IsWritable: true},
			{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
		},
		data,
	)

	// Create and sign transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(wallet.PublicKey()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %v", err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if wallet.PublicKey().Equals(key) {
			return &wallet.PrivateKey
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send transaction
	sig, err := c.client.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return sig.String(), nil
}

// AcceptRequest accepts a pending request for funds
func (c *CashDappClient) AcceptRequest(ctx context.Context, wallet *solana.Wallet, recipient solana.PublicKey) (string, error) {
	// Derive PDA for pending request
	pendingRequestPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(PENDING_REQUEST_SEED),
			recipient.Bytes(),
		},
		c.programID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to derive pending request PDA: %v", err)
	}

	// Get recent blockhash
	recent, err := c.client.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	// Create instruction data
	data := []byte{4, 60, 28, 227, 25, 199, 246, 124} // accept_request discriminator

	// Get pending request account data to get sender and amount
	pendingRequestData, err := c.client.GetAccountInfo(ctx, pendingRequestPDA)
	if err != nil {
		return "", fmt.Errorf("failed to get pending request data: %v", err)
	}

	// Parse sender from pending request data (assuming it's at offset 8)
	sender := solana.PublicKeyFromBytes(pendingRequestData.Value.Data.GetBinary()[8:40])

	// Derive PDAs for cash accounts
	fromCashAccountPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(CASH_ACCOUNT_SEED),
			sender.Bytes(),
		},
		c.programID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to derive sender cash account PDA: %v", err)
	}

	toCashAccountPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(CASH_ACCOUNT_SEED),
			recipient.Bytes(),
		},
		c.programID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to derive recipient cash account PDA: %v", err)
	}

	// Create instruction
	ix := solana.NewInstruction(
		c.programID,
		solana.AccountMetaSlice{
			{PublicKey: pendingRequestPDA, IsSigner: false, IsWritable: true},
			{PublicKey: fromCashAccountPDA, IsSigner: false, IsWritable: true},
			{PublicKey: toCashAccountPDA, IsSigner: false, IsWritable: true},
			{PublicKey: wallet.PublicKey(), IsSigner: true, IsWritable: true},
			{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
		},
		data,
	)

	// Create and sign transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(wallet.PublicKey()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %v", err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if wallet.PublicKey().Equals(key) {
			return &wallet.PrivateKey
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send transaction
	sig, err := c.client.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return sig.String(), nil
}

// DeclineRequest declines a pending request for funds
func (c *CashDappClient) DeclineRequest(ctx context.Context, wallet *solana.Wallet, recipient solana.PublicKey) (string, error) {
	// Derive PDA for pending request
	pendingRequestPDA, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte(PENDING_REQUEST_SEED),
			recipient.Bytes(),
		},
		c.programID,
	)
	if err != nil {
		return "", fmt.Errorf("failed to derive pending request PDA: %v", err)
	}

	// Get recent blockhash
	recent, err := c.client.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	// Create instruction data
	data := []byte{9, 222, 214, 30, 3, 147, 221, 247} // decline_request discriminator

	// Create instruction
	ix := solana.NewInstruction(
		c.programID,
		solana.AccountMetaSlice{
			{PublicKey: pendingRequestPDA, IsSigner: false, IsWritable: true},
			{PublicKey: wallet.PublicKey(), IsSigner: true, IsWritable: true},
			{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
		},
		data,
	)

	// Create and sign transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{ix},
		recent.Value.Blockhash,
		solana.TransactionPayer(wallet.PublicKey()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %v", err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if wallet.PublicKey().Equals(key) {
			return &wallet.PrivateKey
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send transaction
	sig, err := c.client.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return sig.String(), nil
}
