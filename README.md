# auth0-react-go-sample
Auth0 + React + Go sample code

# Auth0でやること
### 0. Auth0 のアカウント作成
### 1. Applications を作成
- Single Page Web Applications
- Settings タブを選択し、以下を入力し、保存
  - Allowed Callback URLs: http://localhost:3000
  - Allowed Logout URLs: http://localhost:3000
  - Allowed Web Origins: http://localhost:3000
- Settings タブの Basic Information の Domain, Client ID を環境変数に設定
### 2. APIs を作成
- Identifier を決める。これを AUDIENCE 環境変数に設定
### 3. Role-Based Access Control (RBAC) を有効化
- APIs -> Settings -> RBAC Settings
  - `Enable RBAC` `Add Permissions in the Access Token` を ON
- APIs -> Permissions
  - Permission (例 `read:messages`) を追加
- User Management -> Roles
  - Role を作成
  - Permissions タブで Permission を追加
  - Users タブで User を追加

# Frontend初回コマンドメモ
```bash
npm create vite@latest . -- --template react-ts
cd react
npm install
npm install @auth0/auth0-react
npm install axios
```

# 参考
- [Hello World Full-Stack Security:
React/JavaScript + Golang Standard Library/Golang](https://developer.auth0.com/resources/code-samples/full-stack/hello-world/basic-role-based-access-control/spa/react-javascript/standard-library-golang)
