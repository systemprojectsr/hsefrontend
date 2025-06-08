from flask import Flask, jsonify
import sqlite3

app = Flask(__name__)

def get_db_connection():
    conn = sqlite3.connect('companies_db.db')
    conn.row_factory = sqlite3.Row
    return conn

def fetch_all_from_table(table_name):
    try:
        conn = get_db_connection()
        cursor = conn.cursor()
        cursor.execute(f"SELECT * FROM {table_name}")
        rows = cursor.fetchall()
        conn.close()
        return [dict(row) for row in rows]
    except sqlite3.Error as e:
        print(f"Error fetching data from table {table_name}: {e}")
        return []

# Получение всех пользователей (/users)
@app.route('/users', methods=['GET'])
def get_all_users():
    users = fetch_all_from_table('users')
    return jsonify(users), 200

# Получение всех компаний (/companies)
@app.route('/companies', methods=['GET'])
def get_all_companies():
    companies = fetch_all_from_table('companies')
    return jsonify(companies), 200

# Получение всех сообщений по связи (/messages/<user_id>)
@app.route('/messages/<int:user_id>', methods=['GET'])
def get_message_for_user(user_id):
    try:
        conn = get_db_connection()
        cursor = conn.cursor()
        cursor.execute("""
            SELECT messages.message_id, messages.message_text
            FROM messages
            JOIN user_companies ON messages.user_company_id = user_companies.user_company_id
            WHERE user_companies.user_id = ?
            """, (user_id,))
        messages = cursor.fetchall()
        conn.close()

        message_list = []
        for message in messages:
            message_list.append({'message_id': message['message_id'], 'message': message['message_text']})

        return jsonify(message_list), 200

    except sqlite3.Error as e:
        print(f"Error getting messages for user {user_id}: {e}")
        return jsonify({'message': 'An error occurred', 'error': str(e)},), 500

# Получение списка компаний для пользователя (/users/<user_id>/companies)
@app.route('/users/<int:user_id>/companies', methods=['GET'])
def get_companies_for_user(user_id):
    try:
        conn = get_db_connection()
        cursor = conn.cursor()
        cursor.execute("""
        SELECT companies.company_id, companies.company_name
        FROM companies
        JOIN user_companies ON companies.company_id = user_companies.company_id
        WHERE user_companies.user_id = ?
        """, (user_id,))
        companies = cursor.fetchall()
        conn.close()

        company_list = []
        for company in companies:
            company_list.append({'company_id': company['company_id'], 'company_name': company['company_name']})

        return jsonify(company_list), 200

    except sqlite3.Error as e:
        print(f"Error getting companies for user {user_id}: {e}")
        return jsonify({'message': 'An error occurred'}), 500

# Получение всех связей (/user_companies)
@app.route('/user_companies', methods=['GET'])
def get_all_user_companies():
    user_companies = fetch_all_from_table('user_companies')
    return jsonify(user_companies), 200

if __name__ == '__main__':
    app.run(debug=False)