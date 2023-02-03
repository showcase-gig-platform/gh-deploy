# gh-deploy

gh-deploy は Github Deployments API を使って Deployment を作成するだけのシンプルなツールです。

Deployment を作成すれば、デプロイのバージョン、日時、実行者の記録が取れるため、プロダクトの変更履歴を確認するのに便利です。

`https://github.com/<YOUR_ORG>/<YOUR_REPO>/deployments`

のような URL から各プロダクトのデプロイ履歴を参照することができます。

Github Actions で Deployment の作成にフックしてデプロイ処理を記述することができます。

https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#deployment

## インストール

```bash
go install github.com/showcase-gig-platform/gh-deploy@latest
```

## 実行方法

```bash
gh-deploy -e <DEPLOY_ENV> -r <YOUR_ORG>/<YOUR_REPO> -t <DEPLOY_TAG>
```

### 例

```bash
gh-deploy -e development -r showcase-gig-platform/your-service -t v0.0.1-dev.1
```
