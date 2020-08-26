# todo

https://gohugo.io/

# FAQ

- O que são containers?

> R: Containers como os que utilizamos hoje surgiram no projeto LXC, fruto de múltiplas funcionalidades que o kernel Linux recebeu nos últimos anos, especialmente na última década, para trabalhar com [Namespaces](https://medium.com/@teddyking/linux-namespaces-850489d3ccf) e [Control Groups (cgroups)](https://github.com/torvalds/linux/blob/master/Documentation/admin-guide/cgroup-v2.rst). Na prática, os containers isolam processos dentro do sistema operacional. Projetos como o Docker adicionam outras funcionalidades aos containers, como, por exemplo, criar uma rede virtual entre containers.

- Por que não usar o `scratch` como stage final do Golang?


- É possível executar container dentro de container?

- É possível criar um container sem utilizar uma imagem de origem?

- É possível usar o sistemd?

- Eu posso colocar várias aplicações no mesmo container?

> R: Sim, mas não deveria. Todo o ecossistema da [Cloud Native Computing Foundation](https://www.cncf.io/) foi criado para trabalhar com o máximo de eficiência e isolamento.

- É melhor utilizar uma imagem da comunidade ou criar uma do zero?

> R: Depende.

- Por que eu preciso me preocupar com o tamanho de uma imagem?

> R: Uma imagem grande é um possível sintoma de que ela possui uma grande e desnecessária quantidade de executáveis, o que pode ser utilizado para explorar vulnerabilidades e realizar ataques. Além disso, podem existir outros problemas, como o código-fonte, senhas e documentos que podem ter sido esquecidos.

- Eu posso colocar testes dentro de um Dockerfile?

> R: Sim, mas não deveria. A utilização do Dockerfile deve ser para exclusivamente instalar todas as dependências e construir uma aplicação. Rotinas de testes devem ser feitas após as imagens serem construídas. Dessa maneira é possível segregar ambientes, um dedicado para construção de aplicações e outro para testar as aplicações.

