CREATE DATABASE acme IF NOT EXISTS;
:use acme
CREATE (banks:Account {accountID: "100", displayName: "Banks"});
CREATE (mainBankAccount:Account {accountID: "101", displayName: "Main bank account"});
