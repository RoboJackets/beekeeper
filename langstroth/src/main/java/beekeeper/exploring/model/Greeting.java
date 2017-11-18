package beekeeper.exploring.model;

import java.math.BigInteger;

/**
 * @author willstuckey
 * @date 3/3/17
 * <p></p>
 */
public class Greeting {
    private BigInteger id;
    private String text;

    public Greeting() {

    }

    public BigInteger getId() {
        return id;
    }

    public void setId(BigInteger id) {
        this.id = id;
    }

    public String getText() {
        return text;
    }

    public void setText(String text) {
        this.text = text;
    }
}
