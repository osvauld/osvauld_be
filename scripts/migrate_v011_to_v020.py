import psycopg2
from psycopg2 import sql

# Connection parameters
dbname = "osvauld"
user = "postgres"
password = "password"
host = "host_url"
# Connect to the PostgreSQL database
conn = psycopg2.connect(dbname=dbname, user=user, password=password, host=host)

# Create a cursor object
cur = conn.cursor()

try:
    # Fetch data from the 'fields' table
    cur.execute("SELECT * FROM fields")
    rows = cur.fetchall()

    # Create Data
    field_groups = {}
    for row in rows:
        field_id, field_name, field_value, field_type, credential_id, user_id, created_at, created_by, updated_at, updated_by = row

        key = f"{field_name}_{credential_id}"
        
        if key not in field_groups:

            field_groups[key] = {
                "field_name": field_name,
                "field_type": field_type,
                "credential_id": credential_id,
                "created_at": created_at,
                "created_by": created_by,
                "updated_at": updated_at,
                "updated_by": updated_by,
                "user_fields": []
            }

        field_groups[key]["user_fields"].append({
            "field_value": field_value,
            "user_id": user_id,
        })

    for key, value in field_groups.items():

        # Insert data into 'field_data'
        cur.execute(
            sql.SQL("INSERT INTO field_data (field_name, field_type, credential_id, created_at, created_by, updated_at, updated_by) VALUES (%s, %s, %s, %s, %s, %s, %s) RETURNING id"),
            (value["field_name"], value["field_type"], value["credential_id"], value["created_at"], value["created_by"], value["updated_at"], value["updated_by"])
        )

        field_id = cur.fetchone()[0]

        for field_value in value["user_fields"]:
            # Insert data into 'field_values'
            cur.execute(
                sql.SQL("INSERT INTO field_values (field_value, user_id, field_id) VALUES (%s, %s, %s)"),
                (field_value["field_value"], field_value["user_id"], field_id)
            )

        # Change admin user type to superadmin
        cur.execute("UPDATE users SET type = 'superadmin' WHERE created_at = (SELECT MIN(created_at) FROM users)")

    # Commit the transaction
    conn.commit()
    print("committed")
    print("Data transferred successfully.")


except Exception as e:
    print(f"An error occurred: {e}")
    conn.rollback()  # Roll back the transaction on error
    raise e

finally:
    # Close the cursor and connection
    cur.close()
    conn.close()


