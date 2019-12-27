#!/usr/bin/env node
import 'source-map-support/register';
import cdk = require('@aws-cdk/core');
import { GundouMireiTypoBotInfraStack } from '../lib/gundou-mirei-typo-bot-infra-stack';

const app = new cdk.App();
new GundouMireiTypoBotInfraStack(app, 'GundouMireiTypoBotInfraStack');
