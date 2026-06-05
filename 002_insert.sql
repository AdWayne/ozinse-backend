-- Базовые роли для системы (взято из макета админки)
INSERT INTO roles (name, permissions) VALUES
('Администратор', '{"all": true}'),
('Пользователь', '{"all": false}')
ON CONFLICT (name) DO NOTHING;

-- Санаттар (Категории)
INSERT INTO categories (name) VALUES
('Телехикая'),
('Мультфильм'),
('Көркем фильм'),
('Деректі фильм'),
('Тв-бағдарлама және реалити-шоу'),
('Ситком'),
('Аниме'),
('Шетел фильмдері')
ON CONFLICT (name) DO NOTHING;

-- Жанрлар (Жанры)
INSERT INTO genres (name) VALUES
('Комедиялар'),
('Отбасымен көретіндер'),
('Ғылыми-танымдық'),
('Ойын-сауық'),
('Ғылыми фантастика және фэнтези'),
('Шытырман оқиғалы'),
('Қысқаметрлі'),
('Музыкалық'),
('Спорттық')
ON CONFLICT (name) DO NOTHING;

-- Жасына сәйкес (Возрастные ограничения)
INSERT INTO age_ratings (range) VALUES
('8-10 жас'),
('10-12 жас'),
('12-14 жас'),
('14-16 жас'),
('16-18 жас')
ON CONFLICT (range) DO NOTHING;



INSERT INTO movies (
    title, 
    description, 
    release_year, 
    director, 
    producer, 
    cover_image_url, 
    category_id, 
    age_rating_id, 
    youtube_video_id
) VALUES (
    'Ғарышқа саяхат', 
    'Адамзаттың ғарышқа алғашқы сапары және жаңа планеталарды зерттеуі туралы қызықты көркем фильм.', 
    2023, 
    'Асқар Үсенов', 
    'Берік Қалиев', 
    'https://example.com/images/space_movie.jpg', 
    (SELECT id FROM categories WHERE name = 'Көркем фильм' LIMIT 1), 
    (SELECT id FROM age_ratings WHERE range = '12-14 жас' LIMIT 1), 
    'dQw4w9WgXcQ'
);


INSERT INTO series (
    title, 
    description, 
    release_year, 
    director, 
    producer, 
    cover_image_url, 
    category_id, 
    age_rating_id
) VALUES (
    'Ауылдастар', 
    'Қаладан ауылға көшіп келген жастардың қызықты да күлкілі оқиғалары туралы ситком.', 
    2022, 
    'Мақсат Оспанов', 
    'Айжан Серікқызы', 
    'https://example.com/images/sitcom_series.jpg', 
    (SELECT id FROM categories WHERE name = 'Ситком' LIMIT 1), 
    (SELECT id FROM age_ratings WHERE range = '14-16 жас' LIMIT 1)
);


INSERT INTO movie_genres (movie_id, genre_id) VALUES (
    (SELECT id FROM movies WHERE title = 'Ғарышқа саяхат' LIMIT 1),
    (SELECT id FROM genres WHERE name = 'Ғылыми фантастика және фэнтези' LIMIT 1)
);

INSERT INTO series_genres (series_id, genre_id) VALUES 
(
    (SELECT id FROM series WHERE title = 'Ауылдастар' LIMIT 1),
    (SELECT id FROM genres WHERE name = 'Комедиялар' LIMIT 1)
),
(
    (SELECT id FROM series WHERE title = 'Ауылдастар' LIMIT 1),
    (SELECT id FROM genres WHERE name = 'Отбасымен көретіндер' LIMIT 1)
);

UPDATE users 
SET role_id = (SELECT id FROM roles WHERE name = 'Администратор') 
WHERE id = 1;