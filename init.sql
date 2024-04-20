CREATE TABLE IF NOT EXISTS income_tax_rates (
    id SERIAL PRIMARY KEY,
    income_level VARCHAR(255),
    min_income DECIMAL(18, 2) NOT NULL,
    max_income DECIMAL(18, 2),
    tax_rate DECIMAL(5, 2) NOT NULL
);


INSERT INTO income_tax_rates (income_level, min_income, max_income, tax_rate)
VALUES
    ('0 - 150,000', 0.00, 150000.00, 0.00),
    ('150,001 - 500,000', 150001.00, 500000.00, 10.00),
    ('500,001 - 1,000,000', 500001.00, 1000000.00, 15.00),
    ('1,000,001 - 2,000,000', 1000001.00, 2000000.00, 20.00),
    ('2,000,001 ขึ้นไป', 2000000.01, 99999999999999, 35.00);


CREATE TABLE IF NOT EXISTS allowances (
    allowance_name VARCHAR(50)  PRIMARY KEY,
    max_allowance DECIMAL(18, 2),
    min_allowance DECIMAL(18, 2)
);

INSERT INTO allowances (allowance_name, max_allowance, min_allowance)
VALUES
    ('Max_Donations_Allowance', 100000.00, NULL),
    ('Default_Personal_Allowance', 60000.00, NULL),
    ('K_Receipt_Max_Allowance', 50000.00, NULL),
    ('Admin_Max_Personal_Allowance', 100000.00, NULL),
    ('Admin_Max_K_Receipt_Allowance', 100000.00, NULL),
    ('Personal_Min_Allowance', NULL, 10000.00);