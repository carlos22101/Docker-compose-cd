
CREATE DATABASE IF NOT EXISTS carlos_DB;
USE carlos_DB;

CREATE TABLE IF NOT EXISTS users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  email VARCHAR(150),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE USER IF NOT EXISTS 'carlos'@'%' IDENTIFIED BY '1234';
GRANT ALL PRIVILEGES ON carlos_DB.* TO 'carlos'@'%';
FLUSH PRIVILEGES;


INSERT INTO users (first_name, last_name, email) VALUES ('Carlos', 'Solis', 'carlos.solis@example.com');
