CREATE TABLE IF NOT EXISTS Users (ID serial PRIMARY KEY, Name VARCHAR (50) NOT NULL, Country VARCHAR (10)  NOT NULL, Points INT NOT NULL);

INSERT INTO Users(Name, Country, Points) VALUES ('Dae', 'tr', 100);
INSERT INTO Users(Name, Country, Points) VALUES ('Bilgitto', 'tr', 200);
INSERT INTO Users(Name, Country, Points) VALUES ('Test_User', 'fr', 300);