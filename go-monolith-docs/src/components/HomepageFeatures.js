import React from 'react';
import clsx from 'clsx';
import styles from './HomepageFeatures.module.css';

const FeatureList = [
  {
    title: 'Easy to Use',
    Svg: require('../../static/img/undraw_docusaurus_mountain.svg').default,
    description: (
      <>
        Go-Monolith has been designed to simplify project development in Go. Almost everything needed for development integrated
        into the Go-Monolith. You shouldn't worry about data migrations, Go-Monolith provides for you flexible migration system.
      </>
    ),
  },
  {
    title: 'Blueprints system',
    Svg: require('../../static/img/undraw_docusaurus_tree.svg').default,
    description: (
      <>
        Go-Monolith fully powered by blueprint methodology. Everything related to one domain system has to be written in one module called
        &nbsp;<a href="https://en.wikipedia.org/wiki/Software_blueprint" target="_blank">blueprint</a>.
      </>
    ),
  },
  {
    title: 'Powered by design patterns',
    Svg: require('../../static/img/undraw_docusaurus_react.svg').default,
    description: (
      <>
          Since Go is a statically typed, compiled programming language, it's not easy to achieve flexibility and easy project expansion.
          For this we have to use design patterns, follow SOLID principles, etc.
      </>
    ),
  },
];

function Feature({Svg, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} alt={title} />
      </div>
      <div className="text--center padding-horiz--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
