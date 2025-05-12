CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(30) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
    Profile-picture VARCHAR(255) '.png''.jpeg''.gif''.jpg',
);

CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content VARCHAR(5000) NOT NULL,
    likes INTEGER DEFAULT 0,
    themes VARCHAR(255),
    created_at VARCHAR(255) DEFAULT CURRENT_TIMESTAMP
    FOREIGN KEY (owner_id) REFERENCES users(id),
);

CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    owner_id INTEGER NOT NULL,
    content NOT NULL,
    created_at VARCHAR(255) DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id),
);

CREATE TABLE identifiant(
    id INTEGER PRIMARY KEY AUOINCREMENT,
    comment_id INTEGER NOT NULL,
    PRIMARY KEY (id, comment_id),
    FOREIGN KEY (id) REFERENCES user(id),
    FOREIGN KEY (comment_id) REFERENCES comments(id),
);