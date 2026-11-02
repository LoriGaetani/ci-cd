
* ogni microservizio ha il suo ciclo di vita isolato. 
* usare il path-filtering per triggerare solo la pipeline del servizio che è stato modificato.
* check pre-merge. 
* si possono aggiungere strumenti per controllare vulnerabilità delle dipendenze, analisi statica del codice, coverage dei test sopra una certa soglia. 
* bisogna allegare ad ogni immagine il SBOM. 
* bisogna cachare le dipendenze usando l'hash dei file con le dipendenze. 
* l'approccio standard è avere un unico branch main e creare branch effimeri per le varie feature e basta.
* è più sicuro che gli step della pipeline di build e test abbiano solo i permessi di lettura e quello che deve fare il push dell'immagine solo quelli di scrittura. conviene fare 2 job (verify e release).
* avendo due job separati dopo che viene accettata la PR viene eseguito il job che fa la release delle immagini su dockerhub. 
* conviene proteggere il branch main per impedire che si possa fare un push diretto.  Settings → Rules → Rulesets → New ruleset → New branch ruleset. selezionando main come target e impostandolo ad Active. 
- Require a pull request before merging: impedisce il normale push diretto su main.
- Require status checks to pass: seleziona i due controlli finali Backend CI e Frontend CI, con GitHub Actions come origine.
- Block force pushes e Restrict deletions: proteggono la storia e l’esistenza del branch.
- Require conversation resolution: utile se usi commenti di revisione sulle PR.
- Require branches to be up to date before merging: consigliabile, ma comporta una nuova esecuzione della CI quando main avanza.
* un'altra pratica da seguire è impostare il controllo di sicurezza sulle dipendenze. 
* mi devo creare un token su docker hub e inserire il mio username come variabile e il token come secret su github. 
* quando uso la semantic-release se non specifico la sintassi non mi crea la versione e quindi non mi pubblica neinte. 

Una best practice è quella vi avere pipeline separate per ogni servizio perchè averne una unica per tutti ha dei difetti: 
* la release è globale: se uno dei servizi fallisce allora non viene rilasciato niente, anche se gli altri servizi sono perfettamente a posto. Invece, ogni servizio dovrebbe essere deployabile indipendentemente. 
* se si usano strumenti come il semantic-release questi danno una versione globale uguale a tutti servizi. ogni servizio invece dovrebbe avere la sua versione indipendente. 

Per garantire che i servizi siano indipendenti devono essere deployabili separatamente, quindi non va bene introdurre anche degli accoppiamenti dal punto di vista temporale: un servizio deve essere rilasciato prima di un altro perchè il suo contratto cambia. La cosa da fare nel caso di cambiamenti nei contratti è che bisogna rendere le modifiche retrocompatibili: il vecchio contratto V1 rimane e aggiungo un V2 con le mie modifiche, i client si aggiornano per agganciarsi a V2, se nessuno usa più V1 e non ci sono problemi si rimuove il vecchio codice. 
Ovviamente se si usano contract test un servizio può rendersi conto se ha rotto qualcosa agli altri e non va in produzione. 





