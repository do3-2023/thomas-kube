-- Create the person table with a UNIQUE constraint on phone_number
CREATE TABLE IF NOT EXISTS person (
    id SERIAL NOT NULL PRIMARY KEY,
    last_name TEXT NOT NULL,
    phone_number TEXT UNIQUE NOT NULL,
    location TEXT NOT NULL
);

-- Insert initial data into the person table
INSERT INTO person (last_name, phone_number, location)
VALUES 
    ('John', '0702030405', 'Marseille'),
    ('Doe', '0603040506', 'Montpellier');
