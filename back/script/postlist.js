const posts = [
    {
      title: "La redstone c'est super cool la team",
      content: "La redstone dans Minecraft est l'un des éléments les plus fascinants et complexes du jeu. Depuis son introduction en 2010, elle a offert aux joueurs un large éventail de possibilités en matière de mécanique de jeu, permettant de créer des systèmes, des machines et des dispositifs de plus en plus sophistiqués. Que ce soit pour réaliser des pièges astucieux, automatiser des tâches répétitives ou même construire des ordinateurs complets, la redstone a transformé Minecraft en un terrain de jeu pour les ingénieurs et les créateurs d'un genre nouveau. La redstone fonctionne de manière similaire à un système électrique. Elle peut être utilisée pour transmettre un signal, mais aussi pour activer des mécanismes et des objets, tels que des portes, des pistons, des lampes, ou même des minecarts. Tout commence par la création de fils de redstone, qui sont la base de tout système de redstone. Ces fils, une fois placés au sol, permettent de transmettre un signal de redstone d’un bloc à un autre. Un signal de redstone peut être activé de différentes manières, mais le plus courant reste l’utilisation d’un levier, d’un bouton ou d’une plaque de pression, qui sert à émettre un impulsion de signal à la redstone. La puissance du signal est importante : elle peut voyager sur une certaine distance, mais plus le signal se propage loin, plus il perd en intensité. Cela signifie que le signal ne pourra pas activer des mécanismes très éloignés à moins d’utiliser des répéteurs. Les répéteurs de redstone sont des blocs qui régénèrent et amplifient le signal, permettant de le faire voyager plus loin. De plus, ils ont l'avantage de pouvoir ajuster la durée du signal et peuvent être orientés pour diriger la redstone dans différentes directions. En plus des fils et des répéteurs, il existe de nombreux autres composants permettant de manipuler et de contrôler les signaux de redstone. Par exemple, les torches de redstone, qui sont une forme de source de signal, peuvent être placées sur des blocs pour fournir un signal constant, et ces torches sont souvent utilisées dans la construction de circuits logiques ou pour activer/désactiver des systèmes. Les pistons et pistons collants sont également des éléments clé de l’ingénierie en redstone. Les pistons peuvent être utilisés pour pousser des blocs lorsqu’un signal de redstone les active, tandis que les pistons collants ont la particularité de tirer les blocs vers eux lorsqu’ils sont activés, ce qui les rend encore plus utiles pour créer des portes secrètes, des mécanismes de transport, ou même des systèmes de tri complexes. Les blocs de commande sont également des éléments puissants pour les utilisateurs avancés. Ils permettent d'exécuter des commandes spécifiques lorsque activés par un signal de redstone. Ces blocs sont extrêmement polyvalents, permettant d’effectuer des actions allant de la téléportation de joueurs à l'exécution de commandes de modification du terrain, en passant par la création de mécanismes totalement automatisés.",
      category: "Redstone"
    },
    {
      title: "Le créatif c'est pour les bouffons",
      content: "La survie dans *Minecraft* est l'une des expériences les plus emblématiques et captivantes du jeu. Dans ce mode, le joueur doit utiliser ses ressources pour survivre dans un monde généré aléatoirement, où il peut rencontrer divers défis allant des créatures hostiles aux conditions climatiques extrêmes. L'objectif principal de la survie est de rester en vie aussi longtemps que possible, tout en développant son environnement, ses compétences et ses ressources. Tout commence par la collecte de ressources de base, comme le bois, les pierres et la nourriture, qui permettent au joueur de se préparer à affronter les dangers qui se dressent devant lui. Au début, le joueur doit rapidement se créer un abri pour se protéger des monstres qui apparaissent la nuit, tels que les zombies, les squelettes et les creepers. L'aspect le plus gratifiant de la survie réside dans la progression du joueur : en collectant des ressources et en explorant, il peut créer de nouveaux outils, des armes plus puissantes, et même des potions pour améliorer ses chances de survie. La construction de refuges et de bases devient essentielle, car le joueur doit se protéger non seulement des monstres, mais aussi de l'environnement. La gestion de la faim et de la santé est cruciale, car les blessures peuvent survenir lors de combats ou d'explorations dangereuses. Les joueurs doivent également s'assurer de cultiver des ressources alimentaires, comme des cultures de blé, de carottes ou des élevages d'animaux, pour maintenir leur niveau de santé. En avançant dans le jeu, le joueur peut explorer de plus en plus loin dans le monde, découvrant des biomes variés, des villages, des temples et d'autres structures générées par le jeu. L'exploration est une facette fascinante de la survie, car elle permet de découvrir de nouvelles ressources, comme le diamant ou l'obsidienne, qui sont nécessaires pour fabriquer des équipements puissants et pour accéder au Nether, une dimension parallèle pleine de dangers. Le Nether est un endroit effrayant et périlleux, où les créatures sont plus fortes et les ressources plus rares, mais il contient également des matériaux essentiels pour avancer dans le jeu, comme le quartz ou les perles de l'Ender, qui sont nécessaires pour accéder à l'End. L'End est la dimension finale du jeu, où se trouve le dragon de l'Ender, un boss très puissant que le joueur doit vaincre pour terminer le jeu. Le mode survie impose des choix stratégiques, car il n'est pas toujours possible de tout faire en même temps. Le joueur doit gérer ses ressources, son équipement et son temps efficacement. Chaque jour dans *Minecraft* en mode survie représente une nouvelle chance de progresser ou de faire face à un nouveau danger. La coopération avec d'autres joueurs en mode multijoueur survie peut également offrir de nouvelles dynamiques, permettant de partager des ressources, de construire ensemble et d'affronter des défis plus grands. La survie, c’est donc une combinaison de stratégie, d'exploration, de gestion des ressources et de combat, qui incite le joueur à s’adapter constamment aux défis du monde de *Minecraft*. C’est un mode qui valorise la patience, la créativité et la persévérance, et qui offre des possibilités infinies de personnalisation de l’expérience. Que ce soit en solo ou en multijoueur, la survie dans *Minecraft* est un défi constant où chaque action compte et où chaque aventure peut mener à une nouvelle découverte.",
      category: "Survie"
    },
  ];
  
  const postsContainer = document.getElementById("posts-container");
  
  posts.forEach(post => {
    const wrapper = document.createElement("div");
    wrapper.className = "post-container";
  
    const postEl = document.createElement("div");
    postEl.className = "post";
  
    postEl.innerHTML = `
      <h2>${post.title}</h2>
      <p>${post.content}</p>
      <div class="category">${post.category}</div>
      <button class="show-more">Voir plus</button>
    `;
  
    const reactions = document.createElement("div");
    reactions.className = "reactions";
    reactions.innerHTML = `
      <div class="reaction-box"><img src="/front/images/like.png" alt="like"></div>
      <div class="reaction-box"><img src="/front/images/comment.png" alt="comment"></div>
      <div class="reaction-box"><img src="/front/images/share.png" alt="share"></div>
    `;
  
    wrapper.appendChild(postEl);
    wrapper.appendChild(reactions);
    postsContainer.appendChild(wrapper);
  
    const showMoreBtn = postEl.querySelector(".show-more");
    showMoreBtn.addEventListener("click", () => {
      postEl.classList.toggle("expanded");
      showMoreBtn.textContent = postEl.classList.contains("expanded") ? "Voir moins" : "Voir plus";
    });
  });
  