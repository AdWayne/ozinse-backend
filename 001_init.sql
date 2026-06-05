DROP TABLE IF EXISTS refresh_tokens CASCADE;
DROP TABLE IF EXISTS featured_content CASCADE;
DROP TABLE IF EXISTS favorite_movies CASCADE;
DROP TABLE IF EXISTS favorite_series CASCADE;
DROP TABLE IF EXISTS episodes CASCADE;
DROP TABLE IF EXISTS seasons CASCADE;
DROP TABLE IF EXISTS movie_genres CASCADE;
DROP TABLE IF EXISTS series_genres CASCADE;
DROP TABLE IF EXISTS genres CASCADE;
DROP TABLE IF EXISTS movies CASCADE;
DROP TABLE IF EXISTS series CASCADE;
DROP TABLE IF EXISTS projects CASCADE; 
DROP TABLE IF EXISTS age_ratings CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles CASCADE;

-- ==========================================
-- 2. СОЗДАНИЕ ТАБЛИЦ (РАЗДЕЛЕНИЕ MOVIES И SERIES)
-- ==========================================

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    permissions JSONB NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(150) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    birth_date DATE,
    role_id INT REFERENCES roles(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE age_ratings (
    id SERIAL PRIMARY KEY,
    range VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    release_year INT,
    director VARCHAR(255),
    producer VARCHAR(255),
    cover_image_url VARCHAR(255),
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    age_rating_id INT REFERENCES age_ratings(id) ON DELETE SET NULL,
    youtube_video_id VARCHAR(50) NOT NULL, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE series (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    release_year INT,
    director VARCHAR(255),
    producer VARCHAR(255),
    cover_image_url VARCHAR(255),
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    age_rating_id INT REFERENCES age_ratings(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE genres (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    icon_url VARCHAR(255)
);

CREATE TABLE movie_genres (
    movie_id INT REFERENCES movies(id) ON DELETE CASCADE,
    genre_id INT REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, genre_id)
);

CREATE TABLE series_genres (
    series_id INT REFERENCES series(id) ON DELETE CASCADE,
    genre_id INT REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (series_id, genre_id)
);

CREATE TABLE seasons (
    id SERIAL PRIMARY KEY,
    series_id INT REFERENCES series(id) ON DELETE CASCADE,
    season_number INT NOT NULL,
    CONSTRAINT unique_series_season UNIQUE (series_id, season_number)
);

CREATE TABLE episodes (
    id SERIAL PRIMARY KEY,
    season_id INT REFERENCES seasons(id) ON DELETE CASCADE,
    episode_number INT NOT NULL,
    title VARCHAR(255),
    youtube_video_id VARCHAR(50) NOT NULL,
    duration INT, 
    CONSTRAINT unique_season_episode UNIQUE (season_id, episode_number)
);

CREATE TABLE favorite_movies (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    movie_id INT REFERENCES movies(id) ON DELETE CASCADE,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, movie_id)
);

CREATE TABLE favorite_series (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    series_id INT REFERENCES series(id) ON DELETE CASCADE,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, series_id)
);

CREATE TABLE featured_content (
    id SERIAL PRIMARY KEY,
    movie_id INT REFERENCES movies(id) ON DELETE CASCADE,
    series_id INT REFERENCES series(id) ON DELETE CASCADE,
    block_type VARCHAR(50) NOT NULL,
    sort_order INT NOT NULL,
    CHECK (
        (movie_id IS NOT NULL AND series_id IS NULL) OR 
        (movie_id IS NULL AND series_id IS NOT NULL)
    ),
    CONSTRAINT unique_block_order UNIQUE (block_type, sort_order)
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT now()
);

