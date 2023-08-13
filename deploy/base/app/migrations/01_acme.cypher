:use neo4j
CREATE DATABASE acme IF NOT EXISTS;
:use acme
CREATE (banks:Account {accountID: "10000", displayName: "Banks"})
CREATE (mainBankAccount:Account {accountID: "10010", displayName: "Main bank account"})
CREATE (oldBankAccount:Account {accountID: "10020", displayName: "Old bank account"})
CREATE (banks)<-[:ChildOf]-(oldBankAccount)
CREATE (banks)<-[:ChildOf]-(mainBankAccount)

CREATE (expenses:Account {accountID: "20000", displayName: "Expenses"})
CREATE (assets:Account {accountID: "20010", displayName: "Assets"})
CREATE (cars:Account {accountID: "20011", displayName: "Cars"})
CREATE (expenses)<-[:ChildOf]-(assets)
CREATE (assets)<-[:ChildOf]-(cars)
;
