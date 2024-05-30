CREATE TABLE block(
    index bigint PRIMARY KEY,
    hash text NOT NULL,
    content json NOT NULL
);

CREATE INDEX block_hash_idx ON block(hash);

CREATE TABLE log(
    address text NOT NULL,
    block bigint NOT NULL,
    content json NOT NULL,
    FOREIGN KEY(block) REFERENCES block(index)
);

CREATE INDEX log_address_idx ON log(address);

CREATE INDEX log_block_idx ON log(block);
