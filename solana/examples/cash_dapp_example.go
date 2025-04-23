package main

import (
	"context"
	"fmt"
	"log"
	"solana/internal/solanaclient"

	"github.com/gagliardetto/solana-go"
)

func main() {
	// 初始化客户端
	client, err := solanaclient.NewCashDappClient("https://api.devnet.solana.com") // 使用 devnet
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// 创建测试钱包
	wallet := solana.NewWallet()
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	ctx := context.Background()

	// 初始化账户
	fmt.Println("Initializing account...")
	sig, err := client.InitializeAccount(ctx, wallet)
	if err != nil {
		log.Fatalf("Failed to initialize account: %v", err)
	}
	fmt.Printf("Account initialized. Transaction signature: %s\n", sig)

	// 存入资金
	amount := uint64(1_000_000_000) // 1 SOL
	fmt.Printf("Depositing %d lamports...\n", amount)
	sig, err = client.DepositFunds(ctx, wallet, amount)
	if err != nil {
		log.Fatalf("Failed to deposit funds: %v", err)
	}
	fmt.Printf("Funds deposited. Transaction signature: %s\n", sig)

	// 创建另一个测试钱包作为朋友
	friendWallet := solana.NewWallet()
	if err != nil {
		log.Fatalf("Failed to create friend wallet: %v", err)
	}

	// 添加朋友
	fmt.Println("Adding friend...")
	sig, err = client.AddFriend(ctx, wallet, friendWallet.PublicKey())
	if err != nil {
		log.Fatalf("Failed to add friend: %v", err)
	}
	fmt.Printf("Friend added. Transaction signature: %s\n", sig)

	// 转账给朋友
	transferAmount := uint64(500_000_000) // 0.5 SOL
	fmt.Printf("Transferring %d lamports to friend...\n", transferAmount)
	sig, err = client.TransferFunds(ctx, wallet, friendWallet.PublicKey(), transferAmount)
	if err != nil {
		log.Fatalf("Failed to transfer funds: %v", err)
	}
	fmt.Printf("Funds transferred. Transaction signature: %s\n", sig)

	// 创建资金请求
	fmt.Println("Creating new pending request...")
	sig, err = client.NewPendingRequest(ctx, friendWallet, wallet.PublicKey(), transferAmount)
	if err != nil {
		log.Fatalf("Failed to create pending request: %v", err)
	}
	fmt.Printf("Pending request created. Transaction signature: %s\n", sig)

	// 接受资金请求
	fmt.Println("Accepting pending request...")
	sig, err = client.AcceptRequest(ctx, wallet, friendWallet.PublicKey())
	if err != nil {
		log.Fatalf("Failed to accept request: %v", err)
	}
	fmt.Printf("Request accepted. Transaction signature: %s\n", sig)

	// 提取资金
	fmt.Printf("Withdrawing %d lamports...\n", transferAmount)
	sig, err = client.WithdrawFunds(ctx, friendWallet, transferAmount)
	if err != nil {
		log.Fatalf("Failed to withdraw funds: %v", err)
	}
	fmt.Printf("Funds withdrawn. Transaction signature: %s\n", sig)
}
