'''

import sqlite3

def create_tables():
    conn = sqlite3.connect('companies_db.db')
    cursor = conn.cursor()

    cursor.execute("""
    CREATE TABLE IF NOT EXISTS messages(
        message_id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_company_id INTEGER NOT NULL,
        message_text TEXT NOT NULL
    )
    """)

    cursor.execute("""
    CREATE TABLE IF NOT EXISTS users (
        user_id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL
    )
    """)

    cursor.execute("""
    CREATE TABLE IF NOT EXISTS companies (
        company_id INTEGER PRIMARY KEY AUTOINCREMENT,
        company_name TEXT NOT NULL
    )
    """)

    cursor.execute("""
    CREATE TABLE IF NOT EXISTS user_companies (
        user_company_id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        company_id INTEGER NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(user_id),
        FOREIGN KEY (company_id) REFERENCES companies(company_id)
    )
    """)
    conn.commit()
    conn.close()

if __name__ == '__main__':
    create_tables()
    print("Tables created successfully!")

'''