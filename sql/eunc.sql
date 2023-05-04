CREATE SCHEMA IF NOT EXISTS eunc;

DROP TABLE IF EXISTS eunc.blocks;
CREATE TABLE eunc.blocks (
    number            BIGINT         NOT NULL,     -- The block number
    hash              VARCHAR(65535) NOT NULL,     -- Hash of the block
    timestamp         BIGINT         NOT NULL,     -- The unix timestamp for when the block was collated
    miner             VARCHAR(65535) NOT NULL,     -- The address of the beneficiary to whom the mining rewards were given
    transaction_count BIGINT         NOT NULL,     -- The number of transactions in the block
    difficulty        NUMERIC(38, 0),              -- Integer of the difficulty for this block
    PRIMARY KEY (number)
);

DROP TABLE IF EXISTS eunc.transactions;
CREATE TABLE eunc.transactions (
  idx               BIGINT         NOT NULL,     -- The tx number
  hash              VARCHAR(65535) NOT NULL,     -- Hash of the transaction
  status            BIGINT         NOT NULL,     -- If transaction send success, 1=success, 0=fail
  block_number      BIGINT         NOT NULL,     -- Block number where this transaction was in
  timestamp         BIGINT         NOT NULL,     -- The unix timestamp for when the block was collated
  value             NUMERIC(38, 0) NOT NULL,     -- Value transferred in Wei
  from_address      VARCHAR(65535) NOT NULL,     -- Address of the sender
  to_address        VARCHAR(65535) DEFAULT NULL, -- Address of the receiver. null when its a contract creation transaction
  gas               BIGINT         NOT NULL,     -- Gas provided by the sender
  gas_price         BIGINT         NOT NULL,     -- Gas price provided by the sender in Wei
  PRIMARY KEY (idx)
);


