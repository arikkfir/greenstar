:use neo4j
CREATE DATABASE acme IF NOT EXISTS;
CALL apoc.util.sleep(1000 * 5);

:use acme
MATCH (n) DETACH DELETE (n);

CREATE (banks:Account {accountID: "10000", displayName: "Banks", icon: "savings"})
CREATE (mainBankAccount:Account {accountID: "10010", displayName: "Main bank account", icon: "savings"})
CREATE (oldBankAccount:Account {accountID: "10020", displayName: "Old bank account", icon: "savings"})
CREATE (banks)<-[:ChildOf]-(oldBankAccount)
CREATE (banks)<-[:ChildOf]-(mainBankAccount)

CREATE (expenses:Account {accountID: "20000", displayName: "Expenses", icon: "shopping_bag"})
CREATE (assets:Account {accountID: "20010", displayName: "Assets", icon: "house"})
CREATE (cars:Account {accountID: "20011", displayName: "Cars", icon: "directions_car"})
CREATE (houses:Account {accountID: "20012", displayName: "Houses", icon: "house"})
CREATE (expenses)<-[:ChildOf]-(assets)
CREATE (assets)<-[:ChildOf]-(cars)
CREATE (assets)<-[:ChildOf]-(houses)

CREATE (mainBankAccount)-[:Transaction {
  txID:"100000",
  date:datetime("2023-08-01T12:34:56.789+0300"),
  referenceID:"abc",
  amount:"123456.789",
  description:"New car!"
}]->(cars)

CREATE (mainBankAccount)-[:Transaction {
  txID:"100010",
  date:datetime("2023-08-02T12:34:56.789+0300"),
  referenceID:"def",
  amount:"234567.890",
  description:"New house!"
}]->(houses)

;
