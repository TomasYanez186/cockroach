// Copyright 2018 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License included
// in the file licenses/BSL.txt and at www.mariadb.com/bsl11.
//
// Change Date: 2022-10-01
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by the Apache License, Version 2.0,
// included in the file licenses/APL.txt and at
// https://www.apache.org/licenses/LICENSE-2.0

import classNames from "classnames";
import React from "react";
import Helmet from "react-helmet";
import { connect } from "react-redux";
import { withRouter, WithRouterProps } from "react-router";

import { doLogin, LoginAPIState } from "src/redux/login";
import { AdminUIState } from "src/redux/state";
import { getDataFromServer } from "src/util/dataFromServer";
import * as docsURL from "src/util/docs";
import { trustIcon } from "src/util/trust";
import InfoBox from "src/views/shared/components/infoBox";

import logo from "assets/crdb.png";
import docsIcon from "!!raw-loader!assets/docs.svg";
import newFeatureIcon from "!!raw-loader!assets/paperPlane.svg";
import "./loginPage.styl";

const version = getDataFromServer().Tag || "UNKNOWN";

interface LoginPageProps {
  loginState: LoginAPIState;
  handleLogin: (username: string, password: string) => Promise<void>;
}

interface LoginPageState {
  username: string;
  password: string;
}

class LoginPage extends React.Component<LoginPageProps & WithRouterProps, LoginPageState> {
  constructor(props: LoginPageProps & WithRouterProps) {
    super(props);
    this.state = {
      username: "",
      password: "",
    };
    // TODO(vilterp): focus username field on mount
  }

  handleUpdateUsername = (evt: React.FormEvent<{ value: string }>) => {
    this.setState({
      username: evt.currentTarget.value,
    });
  }

  handleUpdatePassword = (evt: React.FormEvent<{ value: string }>) => {
    this.setState({
      password: evt.currentTarget.value,
    });
  }

  handleSubmit = (evt: React.FormEvent<any>) => {
    evt.preventDefault();

    this.props.handleLogin(this.state.username, this.state.password)
      .then(() => {
        const { location, router } = this.props;
        if (location.query && location.query.redirectTo) {
          router.push(location.query.redirectTo);
        } else {
          router.push("/");
        }
      });
  }

  renderError() {
    const { error } = this.props.loginState;

    if (!error) {
      return null;
    }

    let message = "Invalid username or password.";
    if (error.message !== "Unauthorized") {
        message = error.message;
    }
    return (
      <div className="login-page__error">Unable to log in: { message }</div>
    );
  }

  render() {
    const inputClasses = classNames("input-text", {
      "input-text--error": !!this.props.loginState.error,
    });

    return (
      <div className="login-page">
        <Helmet>
          <title>Login</title>
        </Helmet>
        <div className="content">
          <section className="section login-page__info">
            <img className="logo" alt="CockroachDB" src={logo} />
            <div className="login-new-feature">
              <div className="login-new-feature-heading">
                <span dangerouslySetInnerHTML={trustIcon(newFeatureIcon)} />
                <span className="login-new-feature-heading__text">New in v2.1</span>
              </div>
              <p className="login-new-feature__blurb">
                A user with a password is required to log in to the UI
                on secure clusters.
              </p>
            </div>
            <InfoBox>
              <h4 className="login-note-box__heading">Note:</h4>
              <p className="login-note-box__blurb">
                Create a user with this SQL command:
              </p>
              <pre className="login-note-box__sql-command">
                <span className="sql-keyword">CREATE USER</span>
                {" "}craig{" "}
                <span className="sql-keyword">WITH PASSWORD</span>
                {" "}
                <span className="sql-string">'cockroach'</span>
                <span className="sql-keyword">;</span>
              </pre>
              <p className="aside">
                <a href={docsURL.adminUILogin} className="login-docs-link" target="_blank">
                  <span className="login-docs-link__icon" dangerouslySetInnerHTML={trustIcon(docsIcon)} />
                  <span className="login-docs-link__text">Read more about configuring login</span>
                </a>
              </p>
            </InfoBox>
          </section>
          <section className="section login-page__form">
            <p className="version">
              Version: <span className="version-tag">{ version }</span>
            </p>
            <div className="form-container">
              <h1 className="heading">Log in to the Web UI</h1>
              {this.renderError()}
              <form onSubmit={this.handleSubmit} className="form-internal" method="post">
                <input
                  type="text"
                  name="username"
                  className={inputClasses}
                  onChange={this.handleUpdateUsername}
                  value={this.state.username}
                  placeholder="Username"
                />
                <input
                  type="password"
                  name="password"
                  className={inputClasses}
                  onChange={this.handleUpdatePassword}
                  value={this.state.password}
                  placeholder="Password"
                />
                <input
                  type="submit"
                  className="submit-button"
                  disabled={this.props.loginState.inProgress}
                  value={this.props.loginState.inProgress ? "Logging in..." : "Log In"}
                />
              </form>
            </div>
          </section>
        </div>
      </div>
    );
  }
}

// tslint:disable-next-line:variable-name
const LoginPageConnected = connect(
  (state: AdminUIState) => {
    return {
      loginState: state.login,
    };
  },
  (dispatch) => ({
    handleLogin: (username: string, password: string) => {
      return dispatch(doLogin(username, password));
    },
  }),
)(withRouter(LoginPage));

export default LoginPageConnected;
