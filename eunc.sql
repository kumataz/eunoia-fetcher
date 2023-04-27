CREATE SCHEMA IF NOT EXISTS eunc;
DROP TABLE IF EXISTS eunc.blocks;

CREATE TABLE eunc.blocks (
    number            BIGINT         NOT NULL,     -- The block number
    hash              VARCHAR(65535) NOT NULL,     -- Hash of the block
    parent_hash       VARCHAR(65535) NOT NULL,     -- Hash of the parent block
    nonce             VARCHAR(65535) NOT NULL,     -- Hash of the generated proof-of-work
    sha3_uncles       VARCHAR(65535) NOT NULL,     -- SHA3 of the uncles data in the block
    logs_bloom        VARCHAR(65535) NOT NULL,     -- The bloom filter for the logs of the block
    transactions_root VARCHAR(65535) NOT NULL,     -- The root of the transaction trie of the block
    state_root        VARCHAR(65535) NOT NULL,     -- The root of the final state trie of the block
    receipts_root     VARCHAR(65535) NOT NULL,     -- The root of the receipts trie of the block
    miner             VARCHAR(65535) NOT NULL,     -- The address of the beneficiary to whom the mining rewards were given
    difficulty        NUMERIC(38, 0),              -- Integer of the difficulty for this block
    -- total_difficulty  NUMERIC(38, 0),              -- Integer of the total difficulty of the chain until this block
    size              BIGINT         NOT NULL,     -- The size of this block in bytes
    extra_data        VARCHAR(65535) DEFAULT NULL, -- The extra data field of this block
    gas_limit         BIGINT         DEFAULT NULL, -- The maximum gas allowed in this block
    gas_used          BIGINT         DEFAULT NULL, -- The total used gas by all transactions in this block
    timestamp         BIGINT         NOT NULL,     -- The unix timestamp for when the block was collated
    transaction_count BIGINT         NOT NULL,     -- The number of transactions in the block
    PRIMARY KEY (number)
);



DROP TABLE IF EXISTS eunc.transactions;

CREATE TABLE eunc.transactions (
  hash              VARCHAR(65535) NOT NULL,     -- Hash of the transaction
  nonce             BIGINT         NOT NULL,     -- The number of transactions made by the sender prior to this one
  block_hash        VARCHAR(65535) NOT NULL,     -- Hash of the block where this transaction was in
  block_number      BIGINT         NOT NULL,     -- Block number where this transaction was in
  transaction_index BIGINT         NOT NULL,     -- Integer of the transactions index position in the block
  from_address      VARCHAR(65535) NOT NULL,     -- Address of the sender
  to_address        VARCHAR(65535) DEFAULT NULL, -- Address of the receiver. null when its a contract creation transaction
  value             NUMERIC(38, 0) NOT NULL,     -- Value transferred in Wei
  gas               BIGINT         NOT NULL,     -- Gas provided by the sender
  gas_price         BIGINT         NOT NULL,     -- Gas price provided by the sender in Wei
  input             VARCHAR(65535) ,             -- The data sent along with the transaction
  input_data        VARCHAR(65535) ,             -- The data(string) sent along with the transaction
  PRIMARY KEY (hash)
);



