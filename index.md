## Workshop Container Best Practices (pt-BR)

Boas práticas para trabalhar com container e ter um ambiente otimizado para melhor performance e segurança. Nosso objetivo será construir uma imagem com o mínimo de layers, menor tamanho e prevenindo os tipos de ataques mais comuns.

Será abordada uma aplicação escrita em Golang, mas os princípios valem para qualquer aplicação.

### 1. Construindo minha aplicação

Código disponível em [steps/01-building-my-app](steps/01-building-my-app)

Vamos começar resolvendo nosso problema:
- Como construir uma aplicação em Golang?

Escolheremos uma imagem popular e que tem muita documentação disponível na internet, a imagem do `ubuntu`:

```dockerfile
FROM ubuntu
# ...
```

Agora instalamos tudo que precisamos e pronto, temos nossa primeira versão. Esse é o passo mais importante para garantirmos que conseguiremos ter nossa aplicação funcional em um container.

Mas existem problemas:

- A imagem não é otimizada para Golang         
- A imagem possui muitos binários              
- A imagem é desnecessariamente grande         
- O usuário padrão é o 'root'                  
- O usuário possui mais privilégios do que o necessário 
- Existem muitas camadas (layers)

Vamos resolver cada um desses problemas.

### 2. Escolhendo uma imagem base melhor

Código disponível em [steps/02-choosing-a-better-image](steps/02-choosing-a-better-image)

O `ubuntu` é uma excelente imagem, mas não possui muitos pacotes que precisamos para trabalhar com Golang. A aplicação que estamos trabalhando é um simples "Hello World" e mesmo assim são necessários instalar centenas de MBs para termos nossa aplicação construída. Dependendo do que fizermos, pode exigir mais pacotes de sistema, o que vai dar ainda mais dor de cabeça. Mesmo para baixar dependências do Golang não é possível, pois o pacote `git` não vem instalado.

A solução é encontrar uma imagem que já possua a maior parte dos pacotes de sistema que precisaremos, e pra isso temos a `golang`. Melhor ainda, nós temos a `golang:1.15`, ou seja, podemos inclusive especificar a versão do Golang que terá todas as ferramentas instaladas para esta versão específica. O mesmo valeria para outras linguagem, que já oferecem uma imagem completa: `node`, `ruby`, `openjdk`, ...

Vamos substituir a imagem `ubuntu` pela `golang:1.15` e fazer as adaptações necessárias. Para isso, sempre leia documentação da imagem que você estiver utilizando:

```dockerfile
FROM golang:1.15
# ...
```

Ótimo! Um problema a menos, agora **a imagem é otimizada para Golang**.

### 3. Multi-stage build

Código disponível em [steps/03-multi-stage-build](steps/03-multi-stage-build)

Agora sim! Temos uma aplicação em Golang com uma imagem com tudo que é necessária para construir nosso "Hello World"! ... E esse é o problema... Ela tem TUDO.

Agora temos uma aplicação que só printa algumas coisas na tela e pesa mais de **850 MB**. Um tamanho absurdo para só escrever uma saída no terminal. Na verdade, é um tamanho desproporcional para a maioria dos microsserviços que você possa precisar fazer no futuro.

O impasse é que a nossa imagem é perfeita para construir uma aplicação Golang, mas não para executá-la. Para essa situação, podemos utilizar **duas ou mais imagens**, que é o que chamamos de [multi-stage build](https://docs.docker.com/develop/develop-images/multistage-build/). As aplicações disso são diversas, a mais comum é contruir uma aplicação ou outro tipo de artefato numa primeira imagem e copiá-la para segunda, e é isso que iremos fazer:

```dockerfile
# Builder
FROM golang:1.15
# ...
RUN go build -o /tmp/main main.go
# ...

# Application
FROM alpine
# ...
COPY --from=0 /tmp/main ./main
# ...
```

Note que existem dois `FROM`, cada um iniciando uma imagem totalmente diferente. Na primeira nós executamos o comando `go build ...` e na segunda copiamos o aplicação `/tmp/main` que foi contruída na primeira.

Agora temos o melhor dos dois mundos, uma imagem completa e com tudo que precisamos para construír nossa aplicação em Golang, e uma imagem enxuta `alpine` apenas com o mínimo para executarmos a aplicação. O resultado é que reduzimos a imagem de mais de **850 MB** para menos de **8 MB**!

Com isso, a nossa **imagem possui poucos binários** e **é pequena**.

### 4. Usuário não-root

Código disponível em [steps/04-non-root-user](steps/04-non-root-user)

Nossa aplicação é um simples "Hello World", mas imagine se ela fosse uma aplicação web pública, disponível para toda a internet. Agora imagine que hackers, ou alguém com uma ferramenta que baixou "por ai", encontraram uma brecha de segurança que permite executar comandos externos, ou seja, praticamente escrever no terminal do seu container. E se o seu usuário do terminal for o `root`...?

Nesse cenário o atacante poderia instalar diversos pacotes dentro do container e ficar minerando criptomoeda na sua infraestrutura. Isso no melhor dos casos. No pior, seria possível orquestrar um ataque para a infraestrutura dentro da própria infraestrutura. Ai se a sua equipe achar que "ah, libera all traffic para toda a VPC, ta dentro da rede mesmo, dá nada", ai sim você estará com muitos problemas.

Para resolver o problema de usarmos o root, basta criarmos outro usuário e torná-lo padrão. Para isso, na segunda imagem, fazemos esse processo:

```dockerfile
# ...
FROM alpine
# ...
RUN adduser -S -D -H -h /app/ madeline
USER madeline
# ...
```

> Alguém já jogou Celeste? ... enfim...

Agora **o usuário padrão não é o `root`** e nossa imagem está muito otimizada e segura.

### 5. Imagem totalmente otimizada

Código disponível em [steps/05-fully-optimized-image](steps/05-fully-optimized-image)

Agora é a hora de ir pra algo mais hardcore.

Eu acredito que se você chegar até o passo 4, já terá construído uma excelente imagem para a maior parte dos cenários. Ainda assim, existem pontos quem podem melhorar ainda mais, inclusive pontos que não têm tanta relação com a imagem, mas sim com a execução do container.

#### Otimização das camadas

Para cada ação no `Dockerfile`, como `RUN` ou `ENV`, é gerada uma nova camada que é referenciada nos metadados da imagem. Talvez mude em versões futuras do Docker, mas existe um "problema" que a cada comando de `RUN` ou que o `RUN` é intercalado com outros comandos como `COPY` ou `ENV`, é gerada uma nova camada como artefato comprimido, ou seja, é como se a cada vez que essa condição ocorre, é gerado um novo "arquivo `.zip`".

Por si só isso não é um problema, mas pode gerar condições indesejadas. Talvez uma das mais comuns é copiar algo para uma camada utilizando `COPY` ou `ADD`, realizar alguma tarefa com `RUN` que resultara em artefatos com dados sensíveis e depois executar um novo `RUN` para apagar esses dados. Na verdade, o que acontece é que o Docker, por usar o [AUFS](https://docs.docker.com/storage/storagedriver/aufs-driver/) utilizará a sobreposição dessas camadas para criar o estado do container onde os arquivos foram excluídos, ou seja, se alguém obter todas as camadas da imagem, conseguirá obter todos os dados sensíveis lá presentes.

#### .dockerignore

Já ouviu falar do `.gitignore`? A função desse arquivo é descrever tudo que deve ser ignorado pelo `git`, evitando que lixo e arquivos sensíveis sejam enviados para um repositório de código. A função do `.dockerignore` é similar: tudo que estiver descrito nele não será enviado para o comando `docker build ...`, mesmo que algo `ADD ./ /app/` esteja descrito num `Dockerfile`.

Isso é importante para evitar que acidentalmente não sejam enviados tanto arquivos como de senhas, mas também arquivos de desenvolvimento local ou até documentação da aplicação. Se a sua aplicação for uma API que trabalhe com transações financeiras, talvez a documentação seja tão crítica quanto o código-fonte.

#### Linux kernel capabilities

Isso não tem tanta relação com a imagem, mas com a execução em si do container. O hospedeiro precisa passar diversos direitos para um container poder executar a aplicação, mas muitos desses direitos podem ser desnecessários para o funcionamento de uma aplicação. Por exemplo, se uma aplicação só executa operações e escreve no `stdout`, não existe razão para ter o direito que abrir uma porta. Nesse cenário, remover essas capacidades do kernel Linux reduzirá a superfície de ataque do seu container.

Como eu disse anteriormente, isso é relacionado à execução do container. Se estiver utilizando o [CLI do Docker](https://docs.docker.com/engine/reference/run/#runtime-privilege-and-linux-capabilities), o comando é `docker run --cap-drop=MKNOD ...`. Variações dessa opção também estão presentes orquestradores de containers como o [Kubernetes](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/) e o [Amazon Elastic Container Service](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definition_parameters.html).

### Fim

Se você chegou até aqui, acredito que tenha conseguido aprender bastante sobre como criar um container seguro e otimizado para suas aplicações.

Agora aproveite uma torta de morango.

### FAQ

- O que são containers?

> R: Containers como os que utilizamos hoje surgiram no projeto LXC, fruto de múltiplas funcionalidades que o kernel Linux recebeu nos últimos anos, especialmente na última década, para trabalhar com [Namespaces](https://medium.com/@teddyking/linux-namespaces-850489d3ccf) e [Control Groups (cgroups)](https://github.com/torvalds/linux/blob/master/Documentation/admin-guide/cgroup-v2.rst). Na prática, os containers isolam processos dentro do sistema operacional. Projetos como o Docker adicionam outras funcionalidades aos containers, como, por exemplo, criar uma rede virtual entre containers.

- Por que não usar o `scratch` como stage final?

> R: Pode utilizar. O problema é que você não possuirá nenhum pacote do sistema operacional, mesmo recursos básicos, como o `sh`. Portanto, se for decidido criar a partir do `scratch`, é fundamental que o desenvolvedor esteja ciente das consequências disso e que a aplicações possua uma cobertura de testes que garanta que o funcionamento da aplicação não será afetada.
> Mais informações: https://docs.docker.com/develop/develop-images/baseimages/

- É possível executar container dentro de container?

> R: Sim. Existem cenários que esse uso se aplica, como por exemplo criar workers do Jenkins em container que conseguem interagir com a API do Docker. O problema de se fazer esse tipo de coisa é que, na prática, é possível hackear praticamente o host inteiro através do container que executa outra container, portanto essa escolha também precisa ser consciente.

- É possível criar um container sem utilizar uma imagem de origem?

> R: Sim. Para isso serve a imagem `scratch`. É a partir dela que imagens como `debian` e `busybox` são construídas. Desde a versão 1.5.0 do Docker, utilizar o `scratch` no `FROM` não cria uma layer extra.

- É possível usar o sistemd?

> Sim, mas não deveria. Containers deveriam 

- Eu posso colocar várias aplicações no mesmo container?

> R: Sim, mas não deveria. Todo o ecossistema da [Cloud Native Computing Foundation](https://www.cncf.io/) foi criado para trabalhar com o máximo de eficiência e isolamento.

- É melhor utilizar uma imagem da comunidade ou criar uma do zero?

> R: Depende. O maior problema é confiar na origem. No caso de imagens como `alpine` e `ubuntu` é confiável pois são as desenvolvedoras das distribuições Linux que construíram e mantêm-nas. Já as outras imagens da comunidade, talvez seja melhor optar pelas imagens que usam o Docker Hub para construí-las, pois é simples auditar o Dockerfile que originou a imagem. Ainda assim, o recomendado é que seja feito um fork do repositório, reconstruída a imagem e colocada em um registry privado da empresa antes de ir para produção.

- Por que eu preciso me preocupar com o tamanho de uma imagem?

> R: Uma imagem grande é um possível sintoma de que ela possui uma grande e desnecessária quantidade de executáveis, o que pode ser utilizado para explorar vulnerabilidades e realizar ataques. Além disso, podem existir outros problemas, como o código-fonte, senhas e documentos que podem ter sido esquecidos.

- Eu posso colocar testes dentro de um Dockerfile?

> R: Sim, mas não deveria. A utilização do Dockerfile deve ser para exclusivamente instalar todas as dependências e construir uma aplicação. Rotinas de testes devem ser feitas após as imagens serem construídas. Dessa maneira é possível segregar ambientes, um dedicado para construção de aplicações e outro para testar as aplicações.
