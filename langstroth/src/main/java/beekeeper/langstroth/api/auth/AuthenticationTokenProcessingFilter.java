package beekeeper.langstroth.api.auth;

import beekeeper.langstroth.model.UserAuthenticationServiceProvider;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.web.filter.GenericFilterBean;

import javax.servlet.FilterChain;
import javax.servlet.ServletException;
import javax.servlet.ServletRequest;
import javax.servlet.ServletResponse;
import javax.servlet.http.HttpServletRequest;
import java.io.IOException;

/**
 * @author willstuckey
 * @date 3/28/17
 * <p></p>
 */
public class AuthenticationTokenProcessingFilter extends GenericFilterBean {

    @Autowired
    private UserAuthenticationServiceProvider userAuthenticationServiceProvider;

    private AuthenticationManager authenticationManager;

    @Override
    public void doFilter(final ServletRequest request,
                         final ServletResponse response,
                         final FilterChain chain)
            throws IOException, ServletException {
        if (request instanceof HttpServletRequest) {
            @SuppressWarnings("unchecked")
            HttpServletRequest httpServletRequest = (HttpServletRequest) request;

            /*
            TODO
            1) process http header for token attribute
            2) validate token
            3) create spring authentication context from user data
            4) set authentication in security context manager
             */
        }

        chain.doFilter(request, response);
    }
}
