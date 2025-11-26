CREATE TABLE product_votes (
       id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
       product_id uuid NOT NULL,
       session_id uuid NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
       machine_id uuid,
       product_name text NOT NULL,
       liked boolean NOT NULL,
       created_at timestamptz NOT NULL DEFAULT now(),
       updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_product_votes_product_id ON product_votes (product_id);
CREATE INDEX idx_product_votes_session_id ON product_votes (session_id);

CREATE UNIQUE INDEX ux_product_votes_session_product
    ON product_votes (session_id, product_id);
