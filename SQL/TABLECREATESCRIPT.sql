DROP TABLE IF EXISTS Product;
DROP TABLE IF EXISTS Master;
DROP TABLE IF EXISTS Client;
DROP TABLE IF EXISTS ProductType;
DROP TABLE IF EXISTS Sizes;
DROP TABLE IF EXISTS Materials;
DROP TABLE IF EXISTS Specializations;
DROP FUNCTION getordersforclientid(integer);

CREATE TABLE Client (
	client_id serial NOT NULL,
	client_fio varchar(30) NOT NULL,
	PRIMARY KEY(client_id)
);

INSERT INTO Client (client_fio)
	VALUES ('Зубенко Михаил Петрович');
INSERT INTO Client (client_fio)
	VALUES ('Волков Евгений Петрович');
INSERT INTO Client (client_fio)
	VALUES ('Цымбалюк Олег Владимирович');
INSERT INTO Client (client_fio)
	VALUES ('Волков Александр Александрович');


CREATE TABLE Specializations (
	specialization_name varchar(30) NOT NULL,
	PRIMARY KEY(specialization_name)
);

INSERT INTO Specializations (specialization_name)
	VALUES ('Одежда');
INSERT INTO Specializations (specialization_name)
	VALUES ('Обувь');
INSERT INTO Specializations (specialization_name)
	VALUES ('Верхняя одежда');
INSERT INTO Specializations (specialization_name)
	VALUES ('Общая специализация');

CREATE TABLE ProductType (
	product_type varchar(30) NOT NULL,
	PRIMARY KEY(product_type)
);

INSERT INTO ProductType (product_type)
	VALUES ('Куртка');
INSERT INTO ProductType (product_type)
	VALUES ('Пальто');
INSERT INTO ProductType (product_type)
	VALUES ('Брюки');
INSERT INTO ProductType (product_type)
	VALUES ('Рубашка');
INSERT INTO ProductType (product_type)
	VALUES ('Платье');
	
CREATE TABLE Master (
	master_id serial NOT NULL,
	master_FIO varchar(30) NOT NULL,
	master_specialization varchar(30) DEFAULT 'Общая специализация',
	PRIMARY KEY(master_id),
	CONSTRAINT fk_specialization
		FOREIGN KEY(master_specialization)
			REFERENCES Specializations(specialization_name)
);

INSERT INTO Master (master_FIO, master_specialization)
	VALUES ('Волкова Татьяна Ивановна','Общая специализация');
INSERT INTO Master (master_FIO, master_specialization)
	VALUES ('Старостин Иван Ильич','Общая специализация');	

CREATE TABLE Sizes (
	s_size varchar(10) NOT NULL,
	PRIMARY KEY(s_size)
);

INSERT INTO Sizes (s_size)
	VALUES ('XS');
INSERT INTO Sizes (s_size)
	VALUES ('S');
INSERT INTO Sizes (s_size)
	VALUES ('M');
INSERT INTO Sizes (s_size)
	VALUES ('L');
INSERT INTO Sizes (s_size)
	VALUES ('XL');

CREATE TABLE Materials (
	m_name varchar(30) NOT NULL,
	PRIMARY KEY(m_name)
);

INSERT INTO Materials (m_name)
	VALUES ('Хлопок');
INSERT INTO Materials (m_name)
	VALUES ('Шелк');
INSERT INTO Materials (m_name)
	VALUES ('Лен');
INSERT INTO Materials (m_name)
	VALUES ('Джинса');

CREATE TABLE Product (
	p_id serial NOT NULL,
	p_type varchar(30) NOT NULL,
	p_cost integer NOT NULL,
	p_size varchar(10) NOT NULL,
	p_material varchar(10) NOT NULL,
	p_master integer NOT NULL,
	p_customer integer NOT NULL,
	PRIMARY KEY(p_id),
	CONSTRAINT fk_material
		FOREIGN KEY(p_material)
			REFERENCES Materials(m_name),
	CONSTRAINT fk_size
		FOREIGN KEY(p_size)
			REFERENCES Sizes(s_size),
	CONSTRAINT fk_type
		FOREIGN KEY(p_type)
			REFERENCES ProductType(product_type),
	CONSTRAINT fk_master
		FOREIGN KEY(p_master)
			REFERENCES Master(master_id),
	CONSTRAINT fk_customer
		FOREIGN KEY(p_customer)
			REFERENCES Client(client_id)
);

INSERT INTO Product (p_type, p_cost, p_size, p_material, p_master, p_customer)
	VALUES ('Брюки', 3000, 'S', 'Джинса', 1, 1);
INSERT INTO Product (p_type, p_cost, p_size, p_material, p_master, p_customer)
	VALUES ('Платье', 3000, 'XS', 'Лен', 2, 1);
INSERT INTO Product (p_type, p_cost, p_size, p_material, p_master, p_customer)
	VALUES ('Рубашка', 1500, 'XL', 'Хлопок', 1, 2);
INSERT INTO Product (p_type, p_cost, p_size, p_material, p_master, p_customer)
	VALUES ('Пальто', 9800, 'M', 'Шелк', 2, 2);

CREATE OR REPLACE FUNCTION getOrdersForClientId(id integer)
	RETURNS TABLE(
		p_id int, p_type varchar(30), p_cost int, p_size varchar(10), p_material varchar(10), master_FIO varchar(30), client_fio varchar(30)
	)
	LANGUAGE plpgsql
AS 
$$
BEGIN
	RETURN QUERY (
	SELECT DISTINCT Product.p_id, Product.p_type, Product.p_cost, Product.p_size, Product.p_material, Master.master_FIO, Client.client_fio
	FROM (Product INNER JOIN Master ON Product.p_master = Master.master_id) INNER JOIN Client ON Product.p_customer = Client.client_id
	WHERE Product.p_customer = id
	);
END;
$$;

SELECT * FROM getOrdersForClientId(1);
