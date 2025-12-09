CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email TEXT NOT NULL UNIQUE,
                       nickname TEXT NOT NULL UNIQUE,
                       password_hash TEXT NOT NULL,

                       first_name TEXT,
                       last_name TEXT,
                       avatar_url TEXT,
                       grade TEXT,
                       major TEXT,
                       city TEXT,
                       description TEXT,

                       created_at TIMESTAMPTZ DEFAULT NOW(),
                       updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE posts (
                       id SERIAL PRIMARY KEY,
                       user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

                       description TEXT,
                       views_count INT DEFAULT 0,

                       created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_posts_user_id ON posts(user_id);

CREATE TABLE files (
                       id SERIAL PRIMARY KEY,
                       post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
                       url TEXT NOT NULL
);

CREATE INDEX idx_files_post_id ON files(post_id);

CREATE TABLE post_likes (
                            id SERIAL PRIMARY KEY,
                            post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
                            user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX uniq_post_likes ON post_likes(post_id, user_id);

CREATE TABLE comments (
                          id SERIAL PRIMARY KEY,

                          post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
                          user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

                          parent_id INT REFERENCES comments(id) ON DELETE CASCADE,

                          text TEXT NOT NULL,
                          created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);

CREATE TABLE comment_likes (
                               id SERIAL PRIMARY KEY,
                               comment_id INT NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
                               user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX uniq_comment_likes ON comment_likes(comment_id, user_id);

CREATE TABLE followers (
                           id SERIAL PRIMARY KEY,
                           user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                           follower_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX uniq_followers ON followers(user_id, follower_id);
CREATE UNIQUE INDEX IF NOT EXISTS uniq_followers ON followers (user_id, follower_id);
CREATE INDEX IF NOT EXISTS idx_followers_user_id ON followers(user_id);
CREATE INDEX IF NOT EXISTS idx_followers_follower_id ON followers(follower_id);
