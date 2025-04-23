package solanaclient

import (
	"context"
	"math"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

type SolanaClient struct {
	client *rpc.Client
}

func NewSolanaClient(RPCEndpoint string) *SolanaClient {
	return &SolanaClient{
		client: rpc.New(RPCEndpoint),
	}
}

func (s *SolanaClient) GetBalance(ctx context.Context, address string) (float64, error) {
	pubKey, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return 0, err
	}

	resp, err := s.client.GetBalance(ctx, pubKey, rpc.CommitmentFinalized)
	if err != nil {
		return 0, err
	}

	// Convert lamports to SOL (1 SOL = 1,000,000,000 lamports)
	balance := float64(resp.Value) / math.Pow10(9)
	return balance, nil
}

func (s *SolanaClient) Transfer(ctx context.Context, fromAddress, toAddress string, amount float64, privateKey string) (string, error) {
	// Convert private key from base58
	// privKeyBytes, err := base58.Decode(privateKey)
	// if err != nil {
	// 	return "", err
	// }

	// Create wallet from private key
	wallet, err := solana.WalletFromPrivateKeyBase58(privateKey)
	if err != nil {
		return "", err
	}

	// Convert amount to lamports
	lamports := uint64(amount * math.Pow10(9))

	// Create transaction
	toPubKey, err := solana.PublicKeyFromBase58(toAddress)
	if err != nil {
		return "", err
	}

	// Get recent blockhash
	recent, err := s.client.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return "", err
	}

	// Create transfer instruction
	transferInstruction := system.NewTransferInstruction(
		lamports,
		wallet.PublicKey(),
		toPubKey,
	).Build()

	// Create and sign transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{transferInstruction},
		recent.Value.Blockhash,
		solana.TransactionPayer(wallet.PublicKey()),
	)
	if err != nil {
		return "", err
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if wallet.PublicKey().Equals(key) {
			return &wallet.PrivateKey
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	// Send transaction
	sig, err := s.client.SendTransaction(ctx, tx)
	if err != nil {
		return "", err
	}

	return sig.String(), nil
}

func (s *SolanaClient) Airdrop(ctx context.Context, address string, amount float64) (string, error) {
	// Only works on devnet/testnet
	pubKey, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return "", err
	}

	// Convert amount to lamports
	lamports := uint64(amount * math.Pow10(9))

	// Request airdrop
	sig, err := s.client.RequestAirdrop(ctx, pubKey, lamports, rpc.CommitmentFinalized)
	if err != nil {
		return "", err
	}

	return sig.String(), nil
}
