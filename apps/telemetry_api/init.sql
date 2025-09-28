-- Create devices table
CREATE TABLE IF NOT EXISTS devices (
    device_id UUID PRIMARY KEY,
    outer_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    device_type VARCHAR(50) NOT NULL,
    home_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create telemetries table
CREATE TABLE IF NOT EXISTS telemetries (
    telemetry_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL,
    value DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_devices_home_id ON devices(home_id);
CREATE INDEX IF NOT EXISTS idx_devices_device_type ON devices(device_type);
CREATE INDEX IF NOT EXISTS idx_telemetries_device_id ON telemetries(device_id);
CREATE INDEX IF NOT EXISTS idx_telemetries_created_at ON telemetries(created_at);
