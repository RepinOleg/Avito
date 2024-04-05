CREATE TABLE feature (
    feature_id SERIAL PRIMARY KEY
);

CREATE TABLE tag (
    tag_id SERIAL PRIMARY KEY
);

CREATE TABLE banner (
    banner_id SERIAL PRIMARY KEY,
    feature_id INT,
    content_title TEXT,
    content_text TEXT,
    content_url TEXT,
    is_active BOOLEAN,
    FOREIGN KEY (feature_id) REFERENCES Feature(feature_id)
);

CREATE TABLE banner_tag (
    banner_id INT,
    tag_id INT,
    PRIMARY KEY (banner_id, tag_id),
    FOREIGN KEY (banner_id) REFERENCES Banner(banner_id),
    FOREIGN KEY (tag_id) REFERENCES Tag(tag_id)
);