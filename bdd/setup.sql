CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(30) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    profile_picture VARCHAR(255)
);

CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content VARCHAR(5000) NOT NULL,
    likes INTEGER DEFAULT 0,
    themes VARCHAR(255),
    created_at VARCHAR(255) DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    owner_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    likes INTEGER DEFAULT 0,
    created_at VARCHAR(255) DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE comment_likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    comment_id INTEGER NOT NULL,
    owner_id INTEGER NOT NULL,
    created_at VARCHAR(255) DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (comment_id) REFERENCES comments(id),
    FOREIGN KEY (owner_id) REFERENCES users(id),
    UNIQUE (comment_id, owner_id)
);

CREATE TABLE post_likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    owner_id INTEGER NOT NULL,
    created_at VARCHAR(255) DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (owner_id) REFERENCES users(id),
    UNIQUE (post_id, owner_id)
);
