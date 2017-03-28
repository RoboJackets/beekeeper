package beekeeper.langstroth.api.auth;

import beekeeper.langstroth.api.APIEndPoints;
import beekeeper.langstroth.model.AuthenticationToken;
import beekeeper.langstroth.model.UserAuthenticationServiceProvider;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

/**
 * @author willstuckey
 * @date 3/28/17
 * <p></p>
 */
@RestController
@RequestMapping(APIEndPoints.Version.V1 + APIEndPoints.Modules.AUTH)
public class AuthenticationEndPoint {

    @Autowired
    private UserAuthenticationServiceProvider userAuthenticationServiceProvider;

    /**
     *
     * @param id
     * @param password
     * @return
     */
    @RequestMapping(method = RequestMethod.GET, value = APIEndPoints.Modules.Auth.WITH_CREDENTIALS)
    private ResponseEntity<AuthenticationToken> getToken(
            @RequestParam("id") final String id,
            @RequestParam("password") final String password) {
        return new ResponseEntity<>(
                userAuthenticationServiceProvider.getToken(id, password),
                HttpStatus.OK);
    }

    /**
     *
     * @param token
     * @return
     */
    @RequestMapping(method = RequestMethod.GET, value = APIEndPoints.Modules.Auth.WITH_TOKEN)
    private ResponseEntity<AuthenticationToken> getToken(
            @RequestParam("token") final String token) {
        return new ResponseEntity<>(
                userAuthenticationServiceProvider.getToken(token),
                HttpStatus.OK);
    }
}
