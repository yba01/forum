CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    username VARCHAR(55) NOT NULL UNIQUE,
    password VARCHAR(150) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS sessions(
    session_id VARCHAR(50) PRIMARY KEY,
    user_id INT,
    username VARCHAR (50) NOT NULL,
    expired_at VARCHAR(50) NOT NULL,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS post (
    id INTEGER PRIMARY KEY,
    user_id INT,
    category_id VARCHAR(255) NOT NULL,
    commentary_id INT,
    title VARCHAR(50) NOT NULL,
    message VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (commentary_id) REFERENCES comments(id)
);

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY,
    user_id INT,
    post_id INT,
    commentary VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES post(id)
);
CREATE TABLE IF NOT EXISTS reaction (
    id_reaction INTEGER PRIMARY KEY,
    post_id INT,
    user_id INT,
    liked INT,
    disliked INT,
    FOREIGN KEY (post_id) REFERENCES post(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS reactioncom (
    id_reactioncom INTEGER PRIMARY KEY,
    com_id INT,
    user_id INT,
    liked INT,
    disliked INT,
    FOREIGN KEY (com_id) REFERENCES comments(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

DROP TABLE IF EXISTS categories;

CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category VARCHAR(50) NOT NULL
);

INSERT INTO categories (category)
VALUES
    ('HTML\CSS'),
    ('Python'),
    ('C'),
    ('JavaScript'),
    ('Go'),
    ('PHP');