CREATE TABLE IF NOT EXISTS income_tax_rates (
    id SERIAL PRIMARY KEY,
    income_level VARCHAR(255) NOT NULL,
    min_income numeric(18, 2) NOT NULL,
    max_income numeric(18, 2) NOT NULL,
    tax_rate numeric(5, 2) NOT NULL
);


INSERT INTO income_tax_rates (income_level, min_income, max_income, tax_rate)
VALUES
    ('0-150,000', 0.00, 150000.00, 0.00),
    ('150,001-500,000', 150001.00, 500000.00, 10.00),
    ('500,001-1,000,000', 500001.00, 1000000.00, 15.00),
    ('1,000,001-2,000,000', 1000001.00, 2000000.00, 20.00),
    ('2,000,001 ขึ้นไป', 2000000.01, 99999999999999, 35.00);


CREATE TABLE IF NOT EXISTS allowances (
	allowance_name varchar(50) PRIMARY KEY NOT NULL,
	max_allowance numeric(18, 2) NOT NULL,
	min_allowance numeric(18, 2) NOT NULL,
	limit_allowance numeric(18, 2) NOT NULL
);



INSERT INTO allowances (allowance_name, max_allowance, min_allowance, limit_allowance)
VALUES('k-receipt', 100000.00, 1.00, 50000.00),
      ('donation', 100000.00, 0, 100000.00),
      ('personal', 100000.00, 10001.00, 60000.00);

