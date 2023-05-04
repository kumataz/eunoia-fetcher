package main

import(
	"fmt"
	"log"
	"time"
	"runtime"
	"context"
	"math/big"
	"flag"
	// "encoding/hex"

	"database/sql"
	_ "github.com/lib/pq"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/core/types"
)

var TxCountDB int64
var BlockCountDB int64
var BlockNumberDB int64
var BlockNumberChain int64

type Config struct{
	EuncRPC string
	DbConn string
	Concurrent int64
	Schema string
	TblBlks string
	TblTxs string
	User string
	Passwd string
	pqInsertTx string
	pqInsertBlock string
}

func OpenDB() (db *sql.DB) {
	// psql postgres://eunoiaad:Trusme123@localhost:5432/eunoiadb
	dbConn := fmt.Sprintf("postgres://%s:%s@%s", cfg.User, cfg.Passwd, cfg.DbConn)
	fmt.Println(dbConn)
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Open DB: %+v\n", db)
	if nil != db{
		fmt.Printf("Open DB: ok\n")
	}

	db.SetMaxOpenConns(int(cfg.Concurrent))
	_, err = db.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s;`, cfg.Schema))
	if nil == err{
		fmt.Printf("Create SCHEMA '%s' ok\n", cfg.Schema)
	} else {
		fmt.Printf("CREATE SCHEMA %s: %+v\n", cfg.Schema, err)
	}

	// deprecated,for better block sync performance
	_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		number            BIGINT         NOT NULL,     -- The block number
		hash              VARCHAR(65535) NOT NULL,     -- Hash of the block
		timestamp         BIGINT         NOT NULL,     -- The unix timestamp for when the block was collated
		miner             VARCHAR(65535) NOT NULL,     -- The address of the beneficiary to whom the mining rewards were given
		transaction_count BIGINT         NOT NULL,     -- The number of transactions in the block
		difficulty        NUMERIC(38, 0) NOT NULL,     -- Integer of the difficulty for this block
		PRIMARY KEY (number)
		);`, cfg.TblBlks),
	)
	if nil != err{
		fmt.Printf("CREATE TABLE %s: %+v\n", cfg.TblBlks, err)
	}
	_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		idx               SERIAL         PRIMARY KEY ,
		hash              VARCHAR(65535) UNIQUE,       -- Hash of the transaction
		block_number      BIGINT         NOT NULL,     -- Block number where this transaction was in
		value             NUMERIC(38, 0) NOT NULL,     -- Value transferred in Wei
		from_address      VARCHAR(65535) NOT NULL,     -- Address of the sender
		to_address        VARCHAR(65535) DEFAULT NULL, -- Address of the receiver. null when its a contract creation transaction
		gas               BIGINT         NOT NULL,     -- Gas provided by the sender
		gas_price         BIGINT         NOT NULL     -- Gas price provided by the sender in Wei
		);`, cfg.TblTxs),
	)
	if nil != err{
		fmt.Printf("CREATE TABLE %s: %+v\n", cfg.TblTxs, err)
	}
	return db
}

func getDbTxCount(db *sql.DB) (n int64){
	rows, err := db.Query(fmt.Sprintf(`SELECT count(value) FROM %s;`, cfg.TblTxs))
	defer rows.Close()
	if nil == err{
		//fmt.Printf("SELECT latest block number %+V \n", rows)
	}
	n = -1
	if nil == rows {
		n = -1
	} else {
		if rows.Next() {
			rows.Scan(&n)
	} else {
		n = -1
	}
	}
	return n //rows
}

func getDbBlockCount(db *sql.DB, lim int64) (n int64){
	var q string
	if (lim <= 0) {
		q = fmt.Sprintf(`SELECT count(number) FROM %s;`, cfg.TblBlks)
	} else {
		q = fmt.Sprintf(`SELECT count(number) FROM %s WHERE numver < %d;`, cfg.TblBlks, lim)
	}
	rows, err := db.Query(q)
	defer rows.Close()
	if nil == err{
	//fmt.Printf("SELECT latest block number %+V \n", rows)
	}
	n = -1
	if nil == rows {
		n = -1
	} else {
		if rows.Next() {
			rows.Scan(&n)
	} else {
		n = -1
		}
	}
	return n //rows
}

func getDbBlockNumber(db *sql.DB) (n int64){
	// fmt.Printf("SELECT max(number) FROM %s;\n", cfg.TblBlks)
	rows, err := db.Query(fmt.Sprintf(`SELECT max(number) FROM %s;`, cfg.TblBlks))
	defer rows.Close()
	if nil == err{
	//fmt.Printf("SELECT latest block number %+V \n", rows)
	}
	n = -1
	if nil == rows {
		n = -1
	} else {
	if rows.Next() {
		rows.Scan(&n)
	} else {
		n = -1
		}
	}
	return n //rows
}

func newClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial(cfg.EuncRPC)
	if err != nil {
		client = nil
	}
	return client, err
}

func getChainBlockNumber(ec *ethclient.Client) (n int64) {
	header, err := ec.HeaderByNumber(context.Background(), nil)
	//Total
	if nil == header || err != nil {
		fmt.Printf("getChainBlockNumber Error: $s\n",err.Error())
		n = -1
		//log.Fatal(err)
	} else {
		n = header.Number.Int64()
	}
	return n
}

const ()


func insertBlock(db *sql.DB, hdr *types.Header, blk *types.Block){
	fmt.Printf("Importing Block %d (%5d TXs) %s\n", blk.NumberU64(), len(blk.Transactions()), blk.Hash().Hex())
	//fmt.Printf(" Block %+v\n", blk)
	_, err := db.Exec(cfg.pqInsertBlock,
		blk.Number().Int64(),    
		blk.Hash().Hex(),
		hdr.Time,
		len(blk.Transactions()),
		hdr.Coinbase.String(),
		hdr.Difficulty.String(),
	)
	if err != nil {
		// panic(err)
		fmt.Printf("Block Error: %s\n", err.Error())
	}
}

func insertTx(db *sql.DB, blk *types.Block, idx int64, tx *types.Transaction){
	fmt.Printf("Importing tx [%5d] %v\n",idx, tx.Hash().Hex())
	msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
	if err != nil {
		//log.Fatal(err)
		fmt.Printf("Skip tx %s, Error: %s\n", tx.Hash().Hex(), err.Error())
		return
	}
	var txto string
	if nil != tx.To() {
		txto = tx.To().Hex()
	} else {
		txto = ""
	}
	_, err = db.Exec(cfg.pqInsertTx,
		idx,
		tx.Hash().Hex(),
		blk.Number().Int64(),
		tx.Value().String(),
		msg.From().Hex(),
		txto,
		tx.Gas(),
		tx.GasPrice().String(),
	)
	if err != nil {
		// panic(err)
		fmt.Printf("Tx Error: %s\n", err.Error())
	}
	//fmt.Printf("%+v\r", r)
}

func fetchWorker (db *sql.DB, ec *ethclient.Client, bn int64, guard chan struct{}){
	lec, err := ethclient.Dial(cfg.EuncRPC) // TODO connection error

	retries := 32
	defer lec.Close()
	hdr,err := lec.HeaderByNumber(context.Background(), big.NewInt(bn))
	for i := 0 ; (nil != err || nil == hdr) && i <  retries ; i++  {
		fmt.Printf("Retries remain: %d, Reading Header[%d] error: %s\n", retries - i, bn, err.Error())
		lec.Close()
		time.Sleep(time.Duration(i*8) * time.Millisecond)
		lec, _ = ethclient.Dial(cfg.EuncRPC)
		hdr, err  = lec.HeaderByNumber(context.Background(), big.NewInt(bn))
	}

	blk, err := lec.BlockByNumber(context.Background(), big.NewInt(bn))
	for i := 0; (nil != err || nil == blk) &&  i < retries ; i++ {
		fmt.Printf("Retries remain: %d, Reading Block[%d] error: %s\n", retries - i, bn, err.Error())
		lec.Close()
		time.Sleep(time.Duration(i*8) * time.Millisecond)
		lec, _ = ethclient.Dial(cfg.EuncRPC)
		blk, err = lec.BlockByNumber(context.Background(), big.NewInt(bn))
	}

	// fmt.Printf("Got header %d, block %d, %s\n", hdr.Number.Int64(), blk.Number().Int64(), blk.Hash().Hex())

	for idx, tx := range(blk.Transactions()) {
		i64 := int64(idx)
		// fmt.Printf(" extracting tx [%5d] %s\n", idx, tx.Hash().Hex())
		// go func(pdb *sql.DB, pblk *types.Block, pidx int64, ptx *types.Transaction) {
		//  insertTx(pdb, pblk, pidx, ptx)
		// }(db, blk, i64, tx)
		insertTx(db, blk, i64, tx)
	}
	insertBlock(db, hdr, blk)
	<-guard
}

func fetcher(db *sql.DB, ec *ethclient.Client, fullySynced bool) {
	runtime.GOMAXPROCS(int(cfg.Concurrent))

	if (BlockCountDB < BlockNumberDB + 1) { // blocks missing
		// TODO find out first Block missing
		// TODO fetch missing blocks
	}

	var begin int64
	// rollback
	if fullySynced {
		begin = BlockNumberDB + 1
	} else if (BlockNumberDB >= cfg.Concurrent){
		begin = BlockNumberDB - cfg.Concurrent
	} else {
		begin = 0
	}

	// deprecated
	r, err := db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE number >= %d;`, cfg.TblBlks, begin))
	fmt.Printf("Rollback: %v\n", r)
	if nil != err {
		fmt.Println(err)
	}

	r, err = db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE block_number >= %d;`, cfg.TblTxs, begin))
	fmt.Printf("Rollback: %v\n", r)
	if nil != err {
		fmt.Println(err)
	}

	guard := make(chan struct{}, (cfg.Concurrent + 7) / 8) // max fetchWorkers
	for bn := begin; bn <= BlockNumberChain; bn++ {
		guard <- struct{}{} // would block if guard channel is already filled
		go fetchWorker (db, ec, bn, guard)
	}
}

func pollDBBlockNumber(db *sql.DB){
	for {
		<-time.After(1 * time.Second)
		// fmt.Println("1 sec")
		n := getDbTxCount(db)
		if TxCountDB != n {
			fmt.Printf("DB Tx count %d, (+%d)\n", n, n - TxCountDB)
			TxCountDB = n
		}
		m := getDbBlockNumber(db)
		if BlockNumberDB != m {
			fmt.Printf("DB blocks %d, (+%d)\n", m, m - BlockNumberDB)
			BlockNumberDB = m
		}
		BlockCountDB = getDbBlockCount(db, 0)
	}
}


func pollChainBlockNumber(ec *ethclient.Client){
	prevBNC := int64(-1)
	for {
		<-time.After(1 * time.Second)
		bnc := getChainBlockNumber(ec)
		if bnc >= 0 {
			BlockNumberChain = bnc
			if prevBNC != BlockNumberChain {
			prevBNC = BlockNumberChain
			fmt.Println("Chain blocks ", BlockNumberChain)
			}
		}
	}
}

var cfg Config = Config {
	EuncRPC: "http://127.0.0.1:8545",
	DbConn: "localhost:5432/eunoiadb",
	Concurrent: 32,
	Schema: "eunc",
	User: "eunoiaad",
	Passwd: "Trusme123",
}

func init() {
	flag.Int64Var     (&cfg.Concurrent, "T", cfg.Concurrent, "Concurrent threads")
	flag.StringVar   (&cfg.Schema,  "s", "eunc", "PostgreSQL database address")
	flag.StringVar   (&cfg.User,    "U", "eunoiaad", "PostgreSQL database address")
	flag.StringVar   (&cfg.Passwd,  "P", "Trusme123", "PostgreSQL database address")
	flag.StringVar   (&cfg.DbConn,  "d", "localhost:5432/eunoiadb", "PostgreSQL database address")
	flag.StringVar   (&cfg.EuncRPC,  "r", "http://127.0.0.1:8545", "RPC URl")
}

// TODO fix table if counnt  does not match last block number 

func main (){
	flag.Parse()
	cfg.TblBlks = cfg.Schema+".blocks"
	cfg.TblTxs = cfg.Schema+".transactions"
	BlockNumberDB = -2
	BlockNumberChain = -2
	TxCountDB = -2
	// total_difficulty deprecated
	cfg.pqInsertBlock = fmt.Sprintf(`INSERT INTO %s (
	number, hash, timestamp, transaction_count, miner, difficulty)
	VALUES ($1, $2, $3, $4, $5, $6);`,
	cfg.TblBlks)

	cfg.pqInsertTx = fmt.Sprintf(`INSERT INTO %s ( 
	idx, hash, block_number, value, from_address, to_address, gas, gas_price)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
	cfg.TblTxs)


	ec, err := newClient()
	if err != nil {
		return
	}
	db := OpenDB()
	go pollDBBlockNumber(db)
	go pollChainBlockNumber(ec)

	time.Sleep(1 * time.Second)
	fullySynced := false
	for ;; {
	time.Sleep(1 * time.Second)
	if -2 != BlockNumberDB && -2 != BlockNumberChain {
		if BlockNumberChain != BlockNumberDB {
			fetcher(db, ec, fullySynced)
		} else {
			if ! fullySynced {
				fmt.Printf("Fully synced, polling new blocks...")
				fullySynced = true
			}
		}
	}
	}
}
