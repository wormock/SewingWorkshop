USE MASTER
GO
DROP DATABASE IF EXISTS SewingWorkShop
GO
CREATE DATABASE SewingWorkShop COLLATE Cyrillic_General_CI_AS
GO
USE SewingWorkShop
GO
DROP TABLE IF EXISTS Product
GO
DROP TABLE IF EXISTS Master
GO
DROP TABLE IF EXISTS Client
GO
DROP TABLE IF EXISTS ProductType
GO
DROP TABLE IF EXISTS Sizes
GO
DROP TABLE IF EXISTS Materials
GO
DROP TABLE IF EXISTS Specializations
GO
CREATE TABLE Client (client_id int IDENTITY (1,1) PRIMARY KEY, client_fio varchar(40) NOT NULL CHECK (client_fio LIKE '[а-яА-Я]%'))
GO
INSERT INTO Client(client_fio) VALUES ('Зубенко Михаил Павлович'), ('Волков Евгений Петрович'), ('Цымбалюк Олег Владимирович'), ('Волков Александр Александрович')
GO
CREATE TABLE Specializations (specialization_name varchar(30) NOT NULL CHECK (specialization_name LIKE '[а-яА-Я]%') PRIMARY KEY)
GO
INSERT INTO Specializations (specialization_name) VALUES ('Одежда'), ('Обувь'), ('Верхняя одежда'), ('Общая специализация')
GO
CREATE TABLE ProductType (product_type varchar(30) NOT NULL CHECK (product_type LIKE '[а-яА-Я]%') PRIMARY KEY)
GO
INSERT INTO ProductType (product_type) VALUES ('Куртка')
GO
INSERT INTO ProductType (product_type) VALUES ('Пальто')
GO
INSERT INTO ProductType (product_type) VALUES ('Брюки')
GO
INSERT INTO ProductType (product_type) VALUES ('Рубашка')
GO
INSERT INTO ProductType (product_type) VALUES ('Платье')
GO
CREATE TABLE Master(master_id int IDENTITY(1,1) PRIMARY KEY, master_FIO varchar(30) NOT NULL CHECK (master_FIO LIKE '[а-яА-Я]%'), master_specialization varchar(30) DEFAULT 'Общая специализация', CONSTRAINT Master_FK_Specialization FOREIGN KEY(master_specialization) REFERENCES Specializations(specialization_name))
GO
INSERT INTO Master(master_FIO, master_specialization) VALUES ('Волкова Татьяна Ивановна','Общая специализация'), ('Старостин Иван Ильич','Общая специализация')
GO
CREATE TABLE Sizes (s_size varchar(10) PRIMARY KEY)
GO 
INSERT INTO Sizes (s_size) VALUES ('XS'), ('S'), ('M'), ('L'), ('XL')
GO
CREATE TABLE Materials (m_name varchar(30) CHECK (m_name LIKE '[а-яА-Я]%') PRIMARY KEY)
INSERT INTO Materials(m_name) VALUES ('Хлопок'),('Шелк'),('Лен'), ('Джинса')
GO
CREATE TABLE Product (p_id int IDENTITY(1,1) PRIMARY KEY, p_type varchar(30) NOT NULL, p_cost int CHECK(p_cost >= 0), p_size varchar(10) NOT NULL, p_material varchar(30) NOT NULL, p_master int NOT NULL, p_customer int NOT NULL, CONSTRAINT fk_size FOREIGN KEY(p_size) REFERENCES Sizes(s_size), CONSTRAINT fk_type FOREIGN KEY(p_type) REFERENCES ProductType(product_type), CONSTRAINT fk_master FOREIGN KEY(p_master) REFERENCES Master(master_id), CONSTRAINT fk_customer FOREIGN KEY(p_customer) REFERENCES Client(client_id))
GO
INSERT INTO Product (p_type, p_cost, p_size, p_material, p_master, p_customer) VALUES ('Брюки', 3000, 'S', 'Джинса', 1, 1), ('Платье', 3000, 'XS', 'Лен', 2, 1), ('Рубашка', 1500, 'XL', 'Хлопок', 1, 2), ('Пальто', 9800, 'M', 'Шелк', 2, 2)
GO
SELECT * FROM Product
GO
SELECT * FROM ProductType
GO
SELECT * FROM Materials
GO
SELECT * FROM Master
GO
SELECT * FROM Client
GO
DROP PROCEDURE IF EXISTS getOrdersForClientId
GO
CREATE PROCEDURE getOrdersForClientId @id int AS SET NOCOUNT ON; SELECT Product.p_id, Product.p_type, Product.p_cost, Product.p_size, Product.p_material, Master.master_FIO, Client.client_fio FROM (Product INNER JOIN Master ON Product.p_master = Master.master_id) INNER JOIN Client ON Product.p_customer = Client.client_id WHERE Product.p_customer = @id
GO
DROP PROCEDURE IF EXISTS addNewMaster
GO
CREATE PROCEDURE addNewMaster @FIO VARCHAR(30), @Specialization VARCHAR(30) AS INSERT INTO Master (master_FIO, master_specialization) VALUES (@FIO, @Specialization)
GO
DROP PROCEDURE IF EXISTS removeMaster 
GO
CREATE PROCEDURE removeMaster @id int AS IF (SELECT COUNT(*) FROM Master WHERE master_id = @id) = 0 BEGIN PRINT('ERROR RECORD WITH THIS PARAMS NOT FOUND') RETURN 2 END IF (SELECT COUNT(*) FROM Product WHERE p_master = @id) != 0 BEGIN PRINT('Unable to delete due to records in linked tables') RETURN 1 END DELETE FROM Master WHERE master_id = @id RETURN 0
GO
DROP PROCEDURE IF EXISTS getMaster
GO
CREATE PROCEDURE getMaster @id int AS SELECT * FROM Master WHERE master_id = @id
GO
DROP PROCEDURE IF EXISTS getLatestWithType
GO
CREATE PROCEDURE getLatestWithType @type VARCHAR(30) AS SELECT DISTINCT Product.p_id, Product.p_type, Product.p_cost, Product.p_size, Product.p_material, Master.master_id, Client.client_id, Master.master_FIO, Client.client_fio  FROM ((Product INNER JOIN Master ON Product.p_master = Master.master_id) INNER JOIN Client ON Product.p_customer = Client.client_id) WHERE p_type = @type
DROP PROCEDURE IF EXISTS editMaster
GO
CREATE PROCEDURE editMaster @ID int, @NEWFIO VARCHAR(30), @NEWSPEC VARCHAR(30) AS UPDATE Master SET master_FIO = @NEWFIO, master_specialization = @NEWSPEC WHERE master_id = @id
GO
DROP PROCEDURE IF EXISTS upCost
GO
CREATE PROCEDURE upCost AS UPDATE Product SET p_cost = p_cost * 1.2
GO
DROP PROCEDURE IF EXISTS getProductCountByCustomer
GO
CREATE PROCEDURE getProductCountByCustomer @ID int AS SELECT COUNT(*) FROM Product WHERE p_customer = @ID