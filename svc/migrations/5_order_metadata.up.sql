CREATE TABLE order_metadata
    (
        id uuid NOT NULL
            CONSTRAINT order_metadata_pk
                PRIMARY KEY,
        order_id uuid NOT NULL,
        "key" text NOT NULL,
        "value" text NOT NULL,
        created_at timestamptz NOT NULL,
        updated_at timestamptz
    );