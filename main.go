package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/dorianneto/gochain/blockchain"
	"github.com/joho/godotenv"
)

var blockchainServer chan []blockchain.Block

func handleConnection(connection net.Conn) {
	defer connection.Close()

	io.WriteString(connection, "Enter a new BPM:")

	scanner := bufio.NewScanner(connection)

	go func() {
		for scanner.Scan() {
			bpm, err := strconv.Atoi(scanner.Text())

			if err != nil {
				log.Printf("%v not a number: %v", scanner.Text(), err)
				continue
			}

			oldBlock := blockchain.Blockchain[len(blockchain.Blockchain)-1]
			newBlock, err := blockchain.GenerateBlock(oldBlock, bpm)

			if err != nil {
				log.Println(err)
			}

			if blockchain.IsBlockValid(newBlock, oldBlock) {
				newBlockchain := append(blockchain.Blockchain, newBlock)
				blockchain.ReplaceChain(newBlockchain)
			}

			blockchainServer <- blockchain.Blockchain
			io.WriteString(connection, "\nEnter a new BPM:")
		}
	}()
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	blockchainServer = make(chan []blockchain.Block)

	timestamp := time.Now()
	genesisBlock := blockchain.Block{Index: 0, Timestamp: timestamp.String(), BPM: 0, Hash: "", PrevHash: ""}

	spew.Dump(genesisBlock)

	blockchain.Blockchain = append(blockchain.Blockchain, genesisBlock)

	server, err := net.Listen("tcp", ":"+os.Getenv("TCP_PORT"))

	if err != nil {
		log.Fatal(err)
	}

	defer server.Close()

	for {
		connection, err := server.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(connection)
	}
}
