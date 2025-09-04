-- Create necessary tables for the API
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20),
    date_of_birth DATE,
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    duration INTEGER NOT NULL,
    genre VARCHAR(100)[],
    rating VARCHAR(10),
    director VARCHAR(100),
    "cast" TEXT[],
    release_date DATE,
    end_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS theaters (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    address TEXT NOT NULL,
    total_halls INTEGER NOT NULL DEFAULT 1,
    contact_phone VARCHAR(20),
    contact_email VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS halls (
    id SERIAL PRIMARY KEY,
    theater_id INTEGER NOT NULL REFERENCES theaters(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    capacity INTEGER NOT NULL,
    screen_type VARCHAR(50),
    has_3d_capability BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS screenings (
    id SERIAL PRIMARY KEY,
    movie_id INTEGER NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    theater_id INTEGER NOT NULL REFERENCES theaters(id) ON DELETE CASCADE,
    hall_id INTEGER NOT NULL REFERENCES halls(id) ON DELETE CASCADE,
    show_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    price_3d DECIMAL(10,2),
    available_seats INTEGER NOT NULL,
    is_3d BOOLEAN DEFAULT FALSE,
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample data for testing
INSERT INTO users (email, password_hash, full_name, phone_number, date_of_birth, email_verified) 
VALUES ('admin@cinema.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Admin User', '08123456789', '1990-01-01', true)
ON CONFLICT (email) DO NOTHING;

INSERT INTO movies (title, description, duration, genre, rating, director, release_date, end_date)
VALUES 
('Avengers: Endgame', 'The epic conclusion to the Infinity Saga', 181, '{"Action","Adventure","Sci-Fi"}', '13+', 'Russo Brothers', '2019-04-26', '2019-07-26'),
('The Batman', 'The Dark Knight investigates corruption in Gotham City', 176, '{"Action","Crime","Drama"}', '13+', 'Matt Reeves', '2022-03-04', '2022-06-04')
ON CONFLICT DO NOTHING;

INSERT INTO theaters (name, address, total_halls, contact_phone, contact_email)
VALUES 
('Cinema XXI Grand Indonesia', 'Jl. M.H. Thamrin No.1, Jakarta', 8, '021-1234567', 'gi@cinema21.com'),
('CGV Pacific Place', 'Jl. Jend. Sudirman Kav. 52-53, Jakarta', 6, '021-7654321', 'pp@cgv.com')
ON CONFLICT DO NOTHING;

INSERT INTO halls (theater_id, name, capacity, screen_type, has_3d_capability)
VALUES 
(1, 'Hall 1', 150, 'IMAX', true),
(1, 'Hall 2', 120, 'Dolby Atmos', true),
(2, 'Studio 1', 100, '4DX', true),
(2, 'Studio 2', 80, 'Regular', false)
ON CONFLICT DO NOTHING;
