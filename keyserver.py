from keyserv import create_app

app = create_app("DevelopmentConfig")

if __name__ == '__main__':
    app.run()
