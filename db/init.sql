
CREATE TABLE  users --para login e cadastro do usuario
(                           --USUARIOS
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY , --sem coordenacao, nao vai ter gargalo quando mudar de DB
    username TEXT UNIQUE NOT NULL,
    senha_hash TEXT  NOT NULL, --as senhas nao precisam ser diferentes
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now()
);

--conta do cliente, carteira
CREATE TABLE accounts
(                                                                  --CONTAS
    user_id UUID PRIMARY KEY REFERENCES users(id),                            --cahve estrangeira ja declarada
    saldo BIGINT   NOT NULL CHECK(saldo>=0),
    reservado BIGINT  NOT NULL DEFAULT 0 CHECK (reservado >= 0 AND reservado <= saldo) --como pego BIGINT? Aqui é em centavos
);

--ACOES
CREATE TABLE assets
(
    --preciso adicionar um ticker que se chama assets, no caso a acao que vai ter silga(chave), nome, setor, preco atual em centavos e por ai vai
    ticker TEXT PRIMARY KEY  NOT NULL, --é a sigla
    nome TEXT UNIQUE  NOT NULL,
    setor TEXT  NOT NULL,
    valor BIGINT CHECK (valor >= 0.0) --acao com centavos, verificar depois
);


--NOME DAS EMPRESAS DOS DOCUMENTOS
INSERT INTO assets VALUES('ARA','aracaju','publico',2090);
INSERT INTO assets VALUES('XLA','Xlag','espacial',4530);
INSERT INTO assets VALUES('SOS','Sosoft','saude', 1020);
INSERT INTO assets VALUES('APR','Appria', 'tecnologia',14520);
INSERT INTO assets VALUES('VOK','Volkjeira','automobilistica',3030);


--orders, as ofertas e intencoes de compra
CREATE TABLE orders
(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,    --É PRIMARY KEY? DEVVE SER COMO SERIAL, VEJA!
    user_id UUID NOT NULL REFERENCES users(id),
    ticker TEXT NOT NULL REFERENCES assets(ticker),
    lado TEXT  NOT NULL CHECK (lado IN ('BUY', 'SELL')), --Ou compra ou vende
    preco BIGINT NOT NULL CHECK (preco>0),                    --0
    quantidade INT NOT NULL CHECK (quantidade>0),
    status TEXT CHECK(status IN('OPEN','PARTIAL','FILLED','CANCELLED','EXPIRED','REJECTED')) DEFAULT 'OPEN',
    criada_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    quantidade_restante INT NOT NULL

    --veja se deve ser direto do accounts ou do users
);



--positions, a carteira
CREATE TABLE positions
(

    user_id UUID NOT NULL REFERENCES users(id),
    ticker TEXT NOT NULL REFERENCES assets(ticker),
    quantidade INT NOT NULL  , --TALVEZ SEJA EM BIGINT, para cripto por exemplo
    reservada INT NOT NULL DEFAULT 0 CHECK(reservada>=0 AND reservada<=quantidade ),         --reservada vai ser o valor ou um estado(sim ou nao)?
    preco_medio BIGINT,
    PRIMARY KEY(user_id,ticker)
);


--trades, cada negociacao executadaa
CREATE TABLE trades
(
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   buy_order_id UUID NOT NULL REFERENCES orders(id),
   sell_order_id UUID NOT NULL REFERENCES orders(id),
   ticker TEXT NOT NULL REFERENCES assets(ticker),
   preco BIGINT NOT NULL CHECK(preco>0),
   quantidade INT NOT NULL CHECK(quantidade>0),
   executado_em  TIMESTAMPTZ NOT NULL DEFAULT now()
);