package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	"SewingWorkshop/pkg/models"

	mssql "github.com/denisenkom/go-mssqldb"
)

// ProductModel - Определяем тип который обертывает пул подключения sql.DB
type ProductModel struct {
	DB *sql.DB
}

func (m *ProductModel) GetProductCountByCustomer(id int) (int, error) {
	stmt := `getProductCountByCustomer @ID = %d`
	row := m.DB.QueryRow(fmt.Sprintf(stmt, id))
	count := &models.ProductCount{}
	err := row.Scan(&count.Count)
	if err != nil {
		return -1, err
	}
	return count.Count, nil
}

func (m *ProductModel) UpdatePrice() error {
	stmt := `upCost`
	_, err := m.DB.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

// Insert - Метод для создания новой заметки в базе дынных.
func (m *ProductModel) Insert(tp string, size string, material string, cost int, master int, customer int) (int, error) {
	// Ниже будет SQL запрос, который мы хотим выполнить. Мы разделили его на две строки
	// для удобства чтения (поэтому он окружен обратными кавычками
	// вместо обычных двойных кавычек).
	stmt := `INSERT INTO Product (p_type, p_cost, p_size, p_material, p_master, p_customer)
	VALUES ($1, $2, $3, $4, $5, $6)`

	// Используем метод Exec() из встроенного пула подключений для выполнения
	// запроса. Первый параметр это сам SQL запрос, за которым следует
	// заголовок заметки, содержимое и срока жизни заметки. Этот
	// метод возвращает объект sql.Result, который содержит некоторые основные
	// данные о том, что произошло после выполнении запроса.
	result, err := m.DB.Exec(stmt, tp, cost, size, material, master, customer)
	if err != nil {
		return 0, err
	}

	// Используем метод LastInsertId(), чтобы получить последний ID
	// созданной записи из таблицу snippets.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Возвращаемый ID имеет тип int64, поэтому мы конвертируем его в тип int
	// перед возвратом из метода.
	return int(id), nil
}

func (m *ProductModel) AddMaster(fio string, specialization string) error {
	stmt := `addNewMaster @FIO = "%s", @Specialization = "%s"`
	_, err := m.DB.Exec(fmt.Sprintf(stmt, fio, specialization))
	if err != nil {
		return err
	}
	return nil
}

func (m *ProductModel) EditMaster(id int, newFio string, newSpec string) error {
	stmt := `editMaster @ID = %d,@NEWFIO = "%s", @NEWSPEC = "%s"`
	_, err := m.DB.Exec((fmt.Sprintf(stmt, id, newFio, newSpec)))
	if err != nil {
		return err
	}
	return nil
}

func (m *ProductModel) GetMaster(id int) (*models.Master, error) {
	stmt := `getMaster @id = %d`
	row := m.DB.QueryRow(fmt.Sprintf(stmt, id))
	master := &models.Master{}
	err := row.Scan(&master.ID, &master.FIO, &master.Specialization)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return master, nil
}

func (m *ProductModel) GetMasters() ([]*models.Master, error) {
	stmt := `SELECT * FROM Master`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	var masters []*models.Master
	for rows.Next() {
		m := &models.Master{}
		err := rows.Scan(&m.ID, &m.FIO, &m.Specialization)
		if err != nil {
			return nil, err
		}
		masters = append(masters, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return masters, nil
}

func (m *ProductModel) DeleteMasterWithId(id int) (int, error) {
	stmt := `removeMaster @id = %d`
	var status mssql.ReturnStatus
	_, err := m.DB.Exec(fmt.Sprintf(stmt, id), &status)
	if err != nil {
		return -1, err
	}
	return int(status), nil
}

func (m *ProductModel) GetOrdersForCustomer(cId int) ([]*models.Product, error) {
	// stmt := `SELECT * FROM getOrdersForClientId($1)`
	stmt := `EXECUTE getOrdersForClientId @id = $1`

	rows, err := m.DB.Query(stmt, cId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product

	for rows.Next() {
		p := &models.Product{}
		err = rows.Scan(&p.ID, &p.Type, &p.Cost, &p.Size, &p.Material, &p.MasterFIO, &p.CustomerFIO)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil

}

// Get - Метод для возвращения данных заметки по её идентификатору ID.
func (m *ProductModel) Get(id int) (*models.Product, error) {
	// SQL запрос для получения данных одной записи.
	stmt := `SELECT * FROM Product WHERE p_id = $1`

	// Используем метод QueryRow() для выполнения SQL запроса,
	// передавая ненадежную переменную id в качестве значения для плейсхолдера
	// Возвращается указатель на объект sql.Row, который содержит данные записи.
	row := m.DB.QueryRow(stmt, id)

	// Инициализируем указатель на новую структуру Product.
	p := &models.Product{}

	// Используйте row.Scan(), чтобы скопировать значения из каждого поля от sql.Row в
	// соответствующее поле в структуре Product. Обратите внимание, что аргументы
	// для row.Scan - это указатели на место, куда требуется скопировать данные
	// и количество аргументов должно быть точно таким же, как количество
	// столбцов в таблице базы данных.
	err := row.Scan(&p.ID, &p.Type, &p.Cost, &p.Size, &p.Material, &p.MasterFIO, &p.CustomerFIO)
	if err != nil {
		// Специально для этого случая, мы проверим при помощи функции errors.Is()
		// если запрос был выполнен с ошибкой. Если ошибка обнаружена, то
		// возвращаем нашу ошибку из модели models.ErrNoRecord.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	// Если все хорошо, возвращается объект Product.
	return p, nil
}

func (m *ProductModel) GetTypes() ([]*models.ProductType, error) {
	stmt := `SELECT * FROM ProductType`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var productTypes []*models.ProductType
	for rows.Next() {
		p := &models.ProductType{}
		err = rows.Scan(&p.TypeName)
		if err != nil {
			return nil, err
		}
		productTypes = append(productTypes, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return productTypes, nil
}

func (m *ProductModel) LatestWithParapms(params string) ([]*models.Product, error) {
	params = fmt.Sprintf("WHERE %s", params)
	stmt := fmt.Sprintf(`SELECT DISTINCT Product.p_id, Product.p_type, Product.p_cost, Product.p_size, Product.p_material, Master.master_id, Client.client_id, Master.master_FIO, Client.client_fio
	FROM (Product INNER JOIN Master ON Product.p_master = Master.master_id) INNER JOIN Client ON Product.p_customer = Client.client_id %s`, params)
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product

	for rows.Next() {
		p := &models.Product{}
		err = rows.Scan(&p.ID, &p.Type, &p.Cost, &p.Size, &p.Material, &p.MasterId, &p.CustomerId, &p.MasterFIO, &p.CustomerFIO)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (m *ProductModel) LatestWithType(pType string) ([]*models.Product, error) {
	stmt := `EXECUTE getLatestWithType @type = $1`

	rows, err := m.DB.Query(stmt, pType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product

	for rows.Next() {
		p := &models.Product{}
		err = rows.Scan(&p.ID, &p.Type, &p.Cost, &p.Size, &p.Material, &p.MasterId, &p.CustomerId, &p.MasterFIO, &p.CustomerFIO)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

// Latest - Метод возвращает последние заказы.
func (m *ProductModel) Latest() ([]*models.Product, error) {
	// Пишем SQL запрос, который мы хотим выполнить.
	// stmt := `SELECT * FROM Product`
	stmt := `SELECT DISTINCT Product.p_id, Product.p_type, Product.p_cost, Product.p_size, Product.p_material, Master.master_id, Client.client_id, Master.master_FIO, Client.client_fio
	FROM (Product INNER JOIN Master ON Product.p_master = Master.master_id) INNER JOIN Client ON Product.p_customer = Client.client_id`

	// Используем метод Query() для выполнения нашего SQL запроса.
	// В ответ мы получим sql.Rows, который содержит результат нашего запроса.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Откладываем вызов rows.Close(), чтобы быть уверенным, что набор результатов из sql.Rows
	// правильно закроется перед вызовом метода Latest(). Этот оператор откладывания
	// должен выполнится *после* проверки на наличие ошибки в методе Query().
	// В противном случае, если Query() вернет ошибку, это приведет к панике
	// так как он попытается закрыть набор результатов у которого значение: nil.
	defer rows.Close()

	// Инициализируем пустой срез для хранения объектов models.Snippets.
	var products []*models.Product

	// Используем rows.Next() для перебора результата. Этот метод предоставляем
	// первый а затем каждую следующею запись из базы данных для обработки
	// методом rows.Scan().
	for rows.Next() {
		// Создаем указатель на новую структуру Product
		p := &models.Product{}
		// Используем rows.Scan(), чтобы скопировать значения полей в структуру.
		// Опять же, аргументы предоставленные в row.Scan()
		// должны быть указателями на место, куда требуется скопировать данные и
		// количество аргументов должно быть точно таким же, как количество
		// столбцов из таблицы базы данных, возвращаемых вашим SQL запросом.
		err = rows.Scan(&p.ID, &p.Type, &p.Cost, &p.Size, &p.Material, &p.MasterId, &p.CustomerId, &p.MasterFIO, &p.CustomerFIO)
		if err != nil {
			return nil, err
		}
		// Добавляем структуру в срез.
		products = append(products, p)
	}

	// Когда цикл rows.Next() завершается, вызываем метод rows.Err(), чтобы узнать
	// если в ходе работы у нас не возникла какая либо ошибка.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Если все в порядке, возвращаем срез с данными.
	return products, nil
}
