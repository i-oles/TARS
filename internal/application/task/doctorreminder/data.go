package doctorreminder

func getMessages(refID string) map[int]string {
	return map[int]string{
		0:  "Chciałbym zapytać, czy zwolnił się może jakiś termin na kolonoskopię do Pana Doktora Witkowskiego.\n\nJestem pacjentem Pana Doktora. Mój numer skierowania to " + refID + ".",
		1:  "Jestem pacjentem Pana Doktora Witkowskiego. Czy mogliby Państwo sprawdzić, czy zwolnił się może jakiś termin na kolonoskopię do Pana Doktora?\n\nNumer mojego skierowania: " + refID + ".",
		2:  "Bardzo proszę o informację, czy pojawił się może wolny termin na kolonoskopię do Pana Doktora Witkowskiego. Jestem pacjentem Pana Doktora.\n\nNumer skierowania: " + refID + ".",
		3:  "Czy mogliby Państwo sprawdzić, czy zwolniło się ostatnio jakieś miejsce na kolonoskopię u Pana Doktora Witkowskiego?\n\nJestem pacjentem Pana Doktora, a numer mojego skierowania to " + refID + ".",
		4:  "Chciałbym uprzejmie zapytać, czy jest obecnie dostępny lub zwolnił się jakiś termin na kolonoskopię do Pana Doktora Witkowskiego.\n\nNumer skierowania: " + refID + ".",
		5:  "Jestem pacjentem Pana Doktora Witkowskiego i chciałbym zapytać, czy pojawił się może wolny termin na kolonoskopię.\n\nMój numer skierowania to " + refID + ".",
		6:  "Czy mogę prosić o sprawdzenie, czy zwolnił się jakiś termin na kolonoskopię u Pana Doktora Witkowskiego?\n\nJestem pacjentem Pana Doktora. Numer skierowania: " + refID + ".",
		7:  "Zwracam się z uprzejmą prośbą o sprawdzenie, czy pojawił się wolny termin na kolonoskopię do Pana Doktora Witkowskiego.\n\nNumer mojego skierowania: " + refID + ".",
		8:  "Czy byłaby możliwość sprawdzenia, czy zwolnił się termin na kolonoskopię do Pana Doktora Witkowskiego?\n\nJestem pacjentem Pana Doktora, numer skierowania " + refID + ".",
		9:  "Bardzo proszę o sprawdzenie dostępności terminu na kolonoskopię u Pana Doktora Witkowskiego. Jestem pacjentem Pana Doktora i zależy mi na wykonaniu badania w tym miejscu.\n\nNumer skierowania: " + refID + ".",
		10: "Czy mogliby Państwo sprawdzić, czy w ostatnim czasie zwolnił się jakiś termin na kolonoskopię do Pana Doktora Witkowskiego?\n\nMój numer skierowania to " + refID + ".",
		11: "Jestem pacjentem Pana Doktora Witkowskiego. Chciałbym zapytać, czy pojawił się może jakiś wolny termin na kolonoskopię.\n\nNumer skierowania: " + refID + ".",
		12: "Uprzejmie proszę o informację, czy zwolnił się może termin na kolonoskopię u Pana Doktora Witkowskiego.\n\nJestem pacjentem Pana Doktora. Numer mojego skierowania to " + refID + ".",
		13: "Chciałbym zapytać, czy jest już może wolny termin kolonoskopii do Pana Doktora Witkowskiego.\n\nNumer skierowania: " + refID + ".",
		14: "Czy mogę prosić o sprawdzenie, czy pojawiło się wolne miejsce na kolonoskopię do Pana Doktora Witkowskiego?\n\nJestem pacjentem Pana Doktora, a numer mojego skierowania to " + refID + ".",
		15: "Bardzo proszę o sprawdzenie, czy nie zwolnił się jakiś termin na kolonoskopię do Pana Doktora Witkowskiego.\n\nJestem pacjentem Pana Doktora. Numer skierowania: " + refID + ".",
		16: "Czy mogliby Państwo sprawdzić, czy aktualnie dostępny jest jakiś termin na kolonoskopię u Pana Doktora Witkowskiego?\n\nJestem pacjentem Pana Doktora, numer mojego skierowania to " + refID + ".",
		17: "Chciałbym uprzejmie zapytać o możliwość wykonania kolonoskopii u Pana Doktora Witkowskiego. Czy zwolnił się może jakiś termin?\n\nNumer skierowania: " + refID + ".",
		18: "Czy jest możliwość sprawdzenia, czy ktoś odwołał wizytę i zwolnił się termin na kolonoskopię do Pana Doktora Witkowskiego?\n\nJestem pacjentem Pana Doktora. Numer skierowania: " + refID + ".",
		19: "Uprzejmie proszę o sprawdzenie, czy pojawił się może wolny termin na kolonoskopię do Pana Doktora Witkowskiego.\n\nJestem pacjentem Pana Doktora, a numer skierowania to " + refID + ".",
	}
}

func getSubjects() map[int]string {
	return map[int]string{
		0: "Termin kolonoskopii – zapytanie",
		1: "Pytanie o termin kolonoskopii",
		2: "Zapytanie o wolny termin kolonoskopii",
		3: "Prośba o sprawdzenie terminu kolonoskopii",
		4: "Pytanie o dostępny termin kolonoskopii",
		5: "Wolny termin na kolonoskopię",
		6: "Zapytanie dotyczące terminu kolonoskopii",
		7: "Prośba o informację w sprawie terminu kolonoskopii",
		8: "Czy zwolnił się termin na kolonoskopię?",
		9: "Dostępność terminu kolonoskopii",
	}
}
