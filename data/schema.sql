-- Ethereum blocks.
CREATE TABLE block(
    index bigint PRIMARY KEY,
    hash bytea NOT NULL,
    content json NOT NULL
);

-- Index to search by block hash.
CREATE INDEX block_hash_idx ON block(hash);

-- Ethereum event log.
CREATE TABLE log(
    id serial PRIMARY KEY, 
    address bytea NOT NULL,
    block bigint,
    content json NOT NULL,
    FOREIGN KEY(block) REFERENCES block(index) ON
    UPDATE
    SET
        NULL -- We'll reuse these rows to avoid VACUUM.
);

-- Index to search by log address.
CREATE INDEX log_address_idx ON log(address);

-- Index to search by log block.
CREATE INDEX log_block_idx ON log(block);
