# SoarSquad
SoarSquad, c'est l'appli parapente qui s'amuse avec vos statistiques, sans aucune pression. Elle fait également office de logbook pour garder une trace de vos vols.
🪂😄 (version 1 seulement le logbook)

# Version 1
- [x] vue d'ensemble
    - [ ] ajout pagination par mois
- [ ] vue d'un vol
- [ ] la progression
- [ ] gestion des exports
- [ ] cron pour recevoir un export annuel



type de fichier pris en compte :
- [x] igc

langue :
- [ ] francais seulement



## stats
- [ ] Temps de vol (par vol, par mois, par année, total, par site)
- [ ] Temps de vol moyen par vol et par année (progression ?)
- [ ] Vitesse max/min
- [ ] Altitude max/min
- [ ] Vario max/min
- [ ] Durée max/min

# Version 2
Les tribus
## Statistiques
### shared

- La fusée ailée (cat 100m, 500m, 1000m) [temps] : La personne qui monte le plus rapidement.
- Marathonien des Nuages [%] : Celui qui a effectué le plus long vol en distance.
- Sorcier du Flottement [%] : Celui qui reste en l'air le plus longtemps pendant un seul vol.
- Navigateur des Gratte-ciels [nombre] : Le parapentiste qui a survolé le plus de zones urbaines.
- Nomade des Sommets [nombre] : La personne ayant visité le plus grand nombre de sites de vol différents.
- Champion du Tournoiement (cat gauche, droite) [nombre] : Le pilote qui effectue le plus grand nombre de tours à gauche ou à droite pendant un vol.
- Virtuose de l'Altitude [nombre] : Le parapentiste qui a atteint l'altitude la plus élevée.
- Gardien de la Nature [nombre] : La personne ayant survolé le plus de zones naturelles protégées.
- Explorateur de l'Aube et du Crépuscule [nombre] : Celui qui réalise le plus de vols tôt le matin ou tard le soir (plus grande amplitude horaire le même jour).
- Pionnier des Hauteurs [nombre] : Celui qui a découvert et volé dans le plus grand nombre de nouveaux sites de vol. (survol site FFVL ?)
- Champion de la Régularité [écart-type ?] : Celui qui maintient la vitesse la plus constante sur le plus grand nombre de vols. (pas de 30s)
- Maître du Retour [nombre] : Celui qui atterrit le plus souvent au point de décollage qu'il a utilisé.
- Sprinteur Céleste [vitesse] : Le parapentiste qui atteint les vitesses les plus élevées en vol.


# Motivations
- Exploration technique : Golang et HTMX (Adieu React ?)
- Besoin personnel : Créer un logbook simple pour suivre ma progression en parapente.
- Une dose d'humour : Intégration de statistiques décalées pour le plaisir.


# Run the code
run
```
./dev.sh
```

if db schema changes run

```
make gen
```
