import os

from keyserv import create_app

if os.environ.get("FLASK_DEBUG"):
    app = create_app("DevelopmentConfig")
else:
    app = create_app("ProductionConfig")

if __name__ == '__main__':
    app.run()
