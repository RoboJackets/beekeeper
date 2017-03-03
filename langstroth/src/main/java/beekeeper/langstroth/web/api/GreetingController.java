package beekeeper.langstroth.web.api;

import beekeeper.langstroth.model.Greeting;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

import java.math.BigInteger;
import java.util.Collection;
import java.util.HashMap;
import java.util.Map;

/**
 * @author willstuckey
 * @date 3/3/17
 * <p></p>
 */
@RestController
public class GreetingController {
    private static BigInteger nextId;
    private static Map<BigInteger, Greeting> greetingMap;

    private static Greeting save(Greeting g) {
        if (greetingMap == null) {
            nextId = BigInteger.ONE;
            greetingMap = new HashMap<>();
        }

        g.setId(nextId);
        nextId = nextId.add(BigInteger.ONE);
        greetingMap.put(g.getId(), g);
        return g;
    }

    static {
        Greeting g1 = new Greeting();
        g1.setText("Hello, World!");
        save(g1);

        Greeting g2 = new Greeting();
        g2.setText("Hola, Mundo!");
        save(g2);

        Greeting g3 = new Greeting();
        g3.setText("Hello, Mundo!");
        save(g3);
    }

    @RequestMapping(value = "/api/v1/greeting",
            method = RequestMethod.GET,
            produces = MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<Collection<Greeting>> getGreetings() {
        Collection<Greeting> retPayload = greetingMap.values();
        return new ResponseEntity<Collection<Greeting>>(retPayload, HttpStatus.OK);
    }
}
