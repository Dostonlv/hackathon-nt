CREATE TYPE tender_status AS ENUM ('open', 'closed', 'awarded');

CREATE TABLE tenders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    deadline TIMESTAMP WITH TIME ZONE NOT NULL,
    budget DECIMAL(15, 2) NOT NULL,
    status tender_status DEFAULT 'open',
    attachment VARCHAR(512),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT budget_positive CHECK (budget > 0)
);

CREATE INDEX idx_tenders_client_id ON tenders(client_id);
CREATE INDEX idx_tenders_status ON tenders(status);

CREATE TABLE bids (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tender_id UUID NOT NULL REFERENCES tenders(id),
    contractor_id UUID NOT NULL REFERENCES users(id),
    price DECIMAL(15, 2) NOT NULL,
    delivery_time INTEGER NOT NULL, -- in days
    comments TEXT,
    status VARCHAR(50) DEFAULT 'open',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT price_positive CHECK (price > 0),
    CONSTRAINT delivery_time_positive CHECK (delivery_time > 0)
);

CREATE INDEX idx_bids_tender_id ON bids(tender_id);
CREATE INDEX idx_bids_contractor_id ON bids(contractor_id);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    message TEXT NOT NULL,
    relation_id UUID,
    type VARCHAR(50) NOT NULL,
    read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
