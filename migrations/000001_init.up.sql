CREATE TABLE IF NOT EXISTS lines
(
    id serial not null unique primary key,
    pair varchar(20) not null,
    val decimal not null,
    description varchar(255) not null,
    typ varchar(20) not null,
    timeframe varchar(20) not null,
);